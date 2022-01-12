package dns

import (
	"math/rand"
	"net"
	"sync"

	"gitlab.com/inetmock/inetmock/internal/app"
	"gitlab.com/inetmock/inetmock/internal/netutils"
)

type RandomIPResolver struct {
	lock   sync.Locker
	Random *rand.Rand
	CIDR   *net.IPNet
}

// nolint:gosec // pseudo-random is desired for this purpose
func NewRandomIPResolver(cidr *net.IPNet) *RandomIPResolver {
	return &RandomIPResolver{
		Random: rand.New(app.RandomSource()),
		CIDR:   cidr,
		lock:   new(sync.Mutex),
	}
}

func (r *RandomIPResolver) Lookup(string) net.IP {
	var (
		ones, bits   = r.CIDR.Mask.Size()
		max          = (1 << (bits - ones)) - 1
		offset, base uint32
	)

	base = netutils.IPToInt32(r.CIDR.IP)
	r.lock.Lock()
	offset = uint32(r.Random.Intn(max))
	r.lock.Unlock()

	return netutils.Uint32ToIP(base + offset)
}
