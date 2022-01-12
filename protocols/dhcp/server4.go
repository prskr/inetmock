package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"

	"gitlab.com/inetmock/inetmock/pkg/logging"
)

const MaxDatagram = 1 << 16

var (
	ErrDropRequest = errors.New("request should be dropped")
	bufPool        = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, MaxDatagram)
			return &buf
		},
	}
)

type Server4 struct {
	conn    *ipv4.PacketConn
	Handler DHCPv4MessageHandler
	Logger  logging.Logger
}

func (s *Server4) Serve(ctx context.Context, conn *ipv4.PacketConn) error {
	s.conn = conn
	if err := s.conn.SetControlMessage(ipv4.FlagInterface, true); err != nil {
		return err
	}
	for ctx.Err() == nil {
		bufBytes := bufPool.Get().(*[]byte)
		b := *bufBytes
		b = b[:MaxDatagram] // Reslice to max capacity in case the buffer in pool was resliced smaller
		n, oob, peer, err := s.conn.ReadFrom(b)
		if err != nil {
			return err
		}
		if msg, err := dhcpv4.FromBytes(b[:n]); err != nil {
			s.Logger.Error("Failed to parse DHCPv4 message", zap.Error(err))
		} else {
			go s.HandleMessage(msg, oob, peer)
		}
		bufPool.Put(bufBytes)
	}
	return ctx.Err()
}

func (s *Server4) HandleMessage(req *dhcpv4.DHCPv4, oob *ipv4.ControlMessage, addr net.Addr) {
	var resp *dhcpv4.DHCPv4
	if r, err := dhcpv4.NewReplyFromRequest(req); err != nil {
		s.Logger.Error("Failed to create response for request", zap.Error(err))
		return
	} else {
		resp = r
	}

	//nolint:exhaustive // other message types are sent by server to the client
	switch mt := req.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		resp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	case dhcpv4.MessageTypeRequest:
		resp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	default:
		s.Logger.Warn("Unhandled message type", zap.String("msg_type", mt.String()))
		return
	}

	if err := s.Handler.Handle(req, resp); err != nil {
		s.Logger.Error("Failed to handle message", zap.Error(err))
		return
	}

	peer, useEthernet := determinePeer(req, resp, addr)

	var woob *ipv4.ControlMessage
	if peer.IP.Equal(net.IPv4bcast) || peer.IP.IsLinkLocalUnicast() || useEthernet {
		// Direct broadcasts, link-local and layer2 unicasts to the interface the request was
		// received on. Other packets should use the normal routing table in
		// case of asymetric routing
		if oob != nil && oob.IfIndex != 0 {
			woob = &ipv4.ControlMessage{IfIndex: oob.IfIndex}
		} else {
			s.Logger.Error("Did not receive interface information")
			return
		}
	}

	if useEthernet {
		intf, err := net.InterfaceByIndex(woob.IfIndex)
		if err != nil {
			s.Logger.Error("Can not get Interface", zap.Error(err))
			return
		}
		err = s.sendEthernet(*intf, resp)
		if err != nil {
			s.Logger.Error("Cannot send Ethernet packet", zap.Error(err))
			return
		}
	} else {
		if _, err := s.conn.WriteTo(resp.ToBytes(), woob, peer); err != nil {
			s.Logger.Error("Failed to write DHCP response", zap.Error(err))
			return
		}
	}
}

// sendEthernet  sends an unicast to the hardware address defined in resp.ClientHWAddr,
// the layer3 destination address is still the broadcast address;
// iface: the interface where the DHCP message should be sent;
// resp: DHCPv4 struct, which should be sent;
func (s *Server4) sendEthernet(iface net.Interface, resp *dhcpv4.DHCPv4) error {
	eth := layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4,
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       resp.ClientHWAddr,
	}
	ip := layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    resp.ServerIPAddr,
		DstIP:    resp.YourIPAddr,
		Protocol: layers.IPProtocolUDP,
		Flags:    layers.IPv4DontFragment,
	}
	udp := layers.UDP{
		SrcPort: dhcpv4.ServerPort,
		DstPort: dhcpv4.ClientPort,
	}

	err := udp.SetNetworkLayerForChecksum(&ip)
	if err != nil {
		return fmt.Errorf("send Ethernet: Couldn't set network layer: %v", err)
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	// Decode a packet
	packet := gopacket.NewPacket(resp.ToBytes(), layers.LayerTypeDHCPv4, gopacket.NoCopy)
	dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4)
	dhcp, ok := dhcpLayer.(gopacket.SerializableLayer)
	if !ok {
		return fmt.Errorf("layer %s is not serializable", dhcpLayer.LayerType().String())
	}
	err = gopacket.SerializeLayers(buf, opts, &eth, &ip, &udp, dhcp)
	if err != nil {
		return fmt.Errorf("cannot serialize layer: %v", err)
	}
	data := buf.Bytes()

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0)
	if err != nil {
		return fmt.Errorf("send Ethernet: Cannot open socket: %v", err)
	}
	defer func() {
		err = syscall.Close(fd)
		if err != nil {
			s.Logger.Error("Send Ethernet: Cannot close socket", zap.Error(err))
		}
	}()

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		s.Logger.Error("Send Ethernet: Cannot set option for socket", zap.Error(err))
	}

	var hwAddr [8]byte
	copy(hwAddr[0:6], resp.ClientHWAddr[0:6])
	ethAddr := syscall.SockaddrLinklayer{
		Protocol: 0,
		Ifindex:  iface.Index,
		Halen:    6,
		Addr:     hwAddr,
	}
	err = syscall.Sendto(fd, data, 0, &ethAddr)
	if err != nil {
		return fmt.Errorf("cannot send frame via socket: %v", err)
	}
	return nil
}

func determinePeer(req, resp *dhcpv4.DHCPv4, pseudoPeer net.Addr) (peer *net.UDPAddr, useEthernet bool) {
	var (
		peerIP     = net.IPv4zero
		clientPort = dhcpv4.ClientPort
	)
	if u, ok := pseudoPeer.(*net.UDPAddr); ok {
		clientPort = u.Port
		peerIP = u.IP
	}
	switch {
	case !req.GatewayIPAddr.IsUnspecified():
		peer = &net.UDPAddr{IP: req.GatewayIPAddr, Port: dhcpv4.ServerPort}
	case resp.MessageType() == dhcpv4.MessageTypeNak:
		peer = &net.UDPAddr{IP: net.IPv4bcast, Port: clientPort}
	case !peerIP.IsUnspecified():
		peer = &net.UDPAddr{IP: peerIP, Port: clientPort}
	case req.IsBroadcast():
		peer = &net.UDPAddr{IP: net.IPv4bcast, Port: clientPort}
	default:
		// sends a layer2 frame so that we can define the destination MAC address
		peer = &net.UDPAddr{IP: resp.YourIPAddr, Port: clientPort}
		useEthernet = true
	}
	return
}
