package netutils

import (
	"errors"
	"net"
)

func IPPortFromAddress(addr net.Addr) (ip net.IP, port int, err error) {
	switch casted := addr.(type) {
	case *net.TCPAddr:
		return casted.IP, casted.Port, nil
	case *net.UDPAddr:
		return casted.IP, casted.Port, nil
	default:
		return nil, 0, errors.New("unknown address type")
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
		defer func() {
			err = errors.Join(err, listener.Close())
		}()
	}

	if addr, ok := listener.Addr().(*net.TCPAddr); !ok {
		return 0, errors.New("not a TCP address")
	} else {
		return addr.Port, nil
	}
}
