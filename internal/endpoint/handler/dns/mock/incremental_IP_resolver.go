package mock

import (
	"net"
	"sync/atomic"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns"
)

type IncrementalIPResolver struct {
	Offset uint32
	CIDR   *net.IPNet
}

func (i *IncrementalIPResolver) Lookup(string) net.IP {
	var (
		ones, bits   = i.CIDR.Mask.Size()
		max          = uint32(1<<(bits-ones)) - 1
		offset, base uint32
	)

	atomic.CompareAndSwapUint32(&i.Offset, max, 0)
	offset = atomic.AddUint32(&i.Offset, 1)

	base = dns.IPToInt32(i.CIDR.IP)
	return dns.Uint32ToIP(base + offset)
}
