package dns

import (
	"math"
	"time"

	mdns "github.com/miekg/dns"
)

func RRHeader(ttl time.Duration, q Question) mdns.RR_Header {
	const minTTLSeconds = 5.0
	ttlSecs := uint32(math.Max(minTTLSeconds, ttl.Seconds()))
	return mdns.RR_Header{
		Name:   q.Name,
		Rrtype: q.Qtype,
		Class:  q.Qclass,
		Ttl:    ttlSecs,
	}
}
