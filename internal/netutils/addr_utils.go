package netutils

import (
	"errors"
	"fmt"
	"net"

	"go.uber.org/multierr"
)

type IPPort struct {
	IP   net.IP
	Port int
}

func (i IPPort) String() string {
	return fmt.Sprintf("%s:%d", i.IP.String(), i.Port)
}

func IPPortFromAddress(addr net.Addr) (ipPort *IPPort, err error) {
	switch casted := addr.(type) {
	case *net.TCPAddr:
		return &IPPort{
			IP:   casted.IP,
			Port: casted.Port,
		}, nil
	case *net.UDPAddr:
		return &IPPort{
			IP:   casted.IP,
			Port: casted.Port,
		}, nil
	default:
		return nil, errors.New("unknown address type")
	}
}

func MustParseMAC(rawMac string) net.HardwareAddr {
	if m, err := net.ParseMAC(rawMac); err != nil {
		panic(err)
	} else {
		return m
	}
}

func RandomPort() (port int, err error) {
	var listener net.Listener
	if l, err := net.Listen("tcp", "127.0.0.1:0"); err != nil {
		return 0, err
	} else {
		listener = l
		defer multierr.AppendInvoke(&err, multierr.Close(listener))
	}

	if addr, ok := listener.Addr().(*net.TCPAddr); !ok {
		return 0, errors.New("not a TCP address")
	} else {
		return addr.Port, nil
	}
}
