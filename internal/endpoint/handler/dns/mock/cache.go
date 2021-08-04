package mock

import (
	"net"
)

type Cache interface {
	PutRecord(host string, address net.IP)
	ForwardLookup(host string) net.IP
	ReverseLookup(address net.IP) (host string, miss bool)
}

type DelegateCache struct {
	OnForwardLookup func(host string) net.IP
	OnReverseLookup func(address net.IP) (host string, miss bool)
	OnPutRecord     func(host string, address net.IP)
}

func (n *DelegateCache) PutRecord(host string, ip net.IP) {
	if n != nil && n.OnPutRecord != nil {
		n.OnPutRecord(host, ip)
	}
}

func (n *DelegateCache) ForwardLookup(host string) net.IP {
	if n != nil && n.OnForwardLookup != nil {
		return n.OnForwardLookup(host)
	}
	return nil
}

func (n *DelegateCache) ReverseLookup(address net.IP) (host string, miss bool) {
	if n != nil && n.OnReverseLookup != nil {
		return n.OnReverseLookup(address)
	}
	return "", true
}
