package dns

import (
	"time"

	mdns "github.com/miekg/dns"
)

func FallbackHandler(handler Handler, resolver IPResolver, ttl time.Duration) Handler {
	return HandlerFunc(func(q Question) (ResourceRecord, error) {
		rr, err := handler.AnswerDNSQuestion(q)
		if err == nil {
			return rr, nil
		}

		ip := resolver.Lookup(q.Name)
		if ip == nil {
			return nil, ErrNoAnswerForQuestion
		}

		switch q.Qtype {
		case mdns.TypeA:
			return &mdns.A{
				A:   ip.To4(),
				Hdr: RRHeader(ttl, q),
			}, nil
		case mdns.TypeAAAA:
			return &mdns.AAAA{
				AAAA: ip.To16(),
				Hdr:  RRHeader(ttl, q),
			}, nil
		default:
			return nil, ErrNoAnswerForQuestion
		}
	})
}
