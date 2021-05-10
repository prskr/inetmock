package netuils

import (
	"errors"
	"fmt"
	"net"
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
