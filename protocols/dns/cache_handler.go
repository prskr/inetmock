package dns

import (
	"math"
	"net"
	"time"

	mdns "github.com/miekg/dns"
)

var _ Handler = (*CacheHandler)(nil)

const (
	minTTLSeconds = 5.0
)

type ResourceRecordCache interface {
	PutRecord(host string, address net.IP)
	ForwardLookup(host string) net.IP
	ReverseLookup(address net.IP) (host string, miss bool)
}

type CacheHandler struct {
	Cache    ResourceRecordCache
	TTL      time.Duration
	Fallback Handler
}

func (h *CacheHandler) AnswerDNSQuestion(q Question) (rr ResourceRecord, err error) {
	switch q.Qtype {
	case mdns.TypeA, mdns.TypeAAAA:
		return h.answerForwardLookup(q)
	case mdns.TypePTR:
		return h.answerReverseLookup(q)
	default:
		return nil, ErrNoAnswerForQuestion
	}
}

func (h CacheHandler) answerForwardLookup(q Question) (rr ResourceRecord, err error) {
	if ip := h.Cache.ForwardLookup(q.Name); ip != nil {
		return &mdns.A{
			A: ip,
			Hdr: mdns.RR_Header{
				Name:   q.Name,
				Class:  mdns.ClassINET,
				Rrtype: q.Qtype,
				Ttl:    h.ttlSeconds(),
			},
		}, nil
	}
	// try to get answer from fallback handler
	if rr, err = h.Fallback.AnswerDNSQuestion(q); err != nil {
		return nil, err
	}

	// put response in cache for further lookups
	switch r := rr.(type) {
	case *mdns.A:
		h.Cache.PutRecord(q.Name, r.A)
	case *mdns.AAAA:
		h.Cache.PutRecord(q.Name, r.AAAA)
	}

	return
}

func (h CacheHandler) answerReverseLookup(q Question) (rr ResourceRecord, err error) {
	ip := ParseInAddrArpa(q.Name)
	if host, miss := h.Cache.ReverseLookup(ip); !miss {
		return &mdns.PTR{
			Ptr: host,
			Hdr: mdns.RR_Header{
				Name:   q.Name,
				Class:  mdns.ClassINET,
				Rrtype: q.Qtype,
				Ttl:    h.ttlSeconds(),
			},
		}, nil
	}

	return nil, ErrNoAnswerForQuestion
}

func (h CacheHandler) ttlSeconds() uint32 {
	return uint32(math.Max(minTTLSeconds, h.TTL.Seconds()))
}
