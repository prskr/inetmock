package mock

import (
	"net"
	"sync"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type IncrementalIPResolver struct {
	lock   sync.Locker
	Offset uint32
	CIDR   *net.IPNet
}

func NewIncrementalIPResolver(cidr *net.IPNet) *IncrementalIPResolver {
	return &IncrementalIPResolver{
		lock: new(sync.Mutex),
		CIDR: cidr,
	}
}

func (i *IncrementalIPResolver) Lookup(string) net.IP {
	i.lock.Lock()
	defer i.lock.Unlock()
	var (
		ones, bits = i.CIDR.Mask.Size()
		max        = uint32(1<<(bits-ones)) - 1
		base       = dns.IPToInt32(i.CIDR.IP)
	)

	if i.Offset >= max {
		i.Offset = 0
	}

	i.Offset += 1

	return dns.Uint32ToIP(base + i.Offset)
}
