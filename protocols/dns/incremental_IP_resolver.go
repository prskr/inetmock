package dns

import (
	"net"
	"sync"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
)

type IncrementalIPResolver struct {
	lock   sync.Mutex
	Offset uint32
	CIDR   *net.IPNet
}

func NewIncrementalIPResolver(cidr *net.IPNet) *IncrementalIPResolver {
	return &IncrementalIPResolver{
		CIDR: cidr,
	}
}

func (i *IncrementalIPResolver) Lookup(string) net.IP {
	i.lock.Lock()
	defer i.lock.Unlock()
	var (
		ones, bits = i.CIDR.Mask.Size()
		max        = uint32(1<<(bits-ones)) - 1
		base       = netutils.IPToInt32(i.CIDR.IP)
	)

	if i.Offset >= max {
		i.Offset = 0
	}

	i.Offset += 1

	return netutils.Uint32ToIP(base + i.Offset)
}
