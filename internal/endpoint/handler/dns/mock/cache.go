package mock

import (
	"net"
)

type Cache interface {
	PutRecord(host string, address net.IP)
	ForwardLookup(host string) net.IP
	ReverseLookup(address net.IP) (host string, miss bool)
}

type NoOpCache struct {
}

func (n NoOpCache) PutRecord(string, net.IP) {
}

func (n NoOpCache) ForwardLookup(string) net.IP {
	return nil
}

func (n NoOpCache) ReverseLookup(net.IP) (host string, miss bool) {
	return "", true
}
