package dhcp

import (
	"encoding/base64"
	"errors"
	"hash/fnv"
	"net"
	"path"
	"sync"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/multierr"

	"inetmock.icb4dc0.de/inetmock/internal/netutils"
	"inetmock.icb4dc0.de/inetmock/internal/state"
)

const rangeHandlerStatePrefix = "range"

type RangeLease struct {
	IP  net.IP
	MAC net.HardwareAddr
}

type RangeMessageHandler struct {
	lock     sync.Mutex
	rangeKey string
	Store    state.KVStore
	TTL      time.Duration
	StartIP  net.IP
	EndIP    net.IP
}

func (h *RangeMessageHandler) Handle(req, resp *dhcpv4.DHCPv4) (err error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	lease := new(RangeLease)
	defer func() {
		if err == nil {
			resp.YourIPAddr = lease.IP
			resp.Options.Update(dhcpv4.OptIPAddressLeaseTime(h.TTL))
		}
	}()

	if h.rangeKey == "" {
		h.calcRangeKey()
	}

	macKey := path.Join(rangeHandlerStatePrefix, h.rangeKey, req.ClientHWAddr.String())
	return h.Store.ReadWriteTransaction(func(rw state.TxnReaderWriter) error {
		if err := rw.Get(macKey, lease); err == nil {
			ipKey := path.Join(rangeHandlerStatePrefix, h.rangeKey, lease.IP.String())
			return multierr.Combine(
				rw.Set(macKey, lease, state.WithTTL(h.TTL)),
				rw.Set(ipKey, lease, state.WithTTL(h.TTL)),
			)
		}

		var leases []RangeLease
		if err := rw.GetAll(path.Join(rangeHandlerStatePrefix, h.rangeKey), &leases); err != nil {
			return err
		}
		lookup := leasesToLookup(leases)
		endIPVal := netutils.IPToInt32(h.EndIP)
		for ipVal := netutils.IPToInt32(h.StartIP); ipVal < endIPVal; ipVal++ {
			if _, ok := lookup[ipVal]; !ok {
				lease.MAC = req.ClientHWAddr
				lease.IP = netutils.Uint32ToIP(ipVal)
				ipKey := path.Join(rangeHandlerStatePrefix, h.rangeKey, lease.IP.String())
				return multierr.Combine(
					rw.Set(macKey, lease, state.WithTTL(h.TTL)),
					rw.Set(ipKey, lease, state.WithTTL(h.TTL)),
				)
			}
		}

		return errors.New("no free IP in range")
	})
}

func (h *RangeMessageHandler) calcRangeKey() {
	hash := fnv.New32a()
	h.rangeKey = base64.URLEncoding.EncodeToString(hash.Sum(append(h.StartIP, h.EndIP...)))
}

func leasesToLookup(leases []RangeLease) map[uint32]RangeLease {
	res := make(map[uint32]RangeLease)
	for idx := range leases {
		res[netutils.IPToInt32(leases[idx].IP)] = leases[idx]
	}
	return res
}
