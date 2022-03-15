package dns_test

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	dnsmock "gitlab.com/inetmock/inetmock/internal/mock/dns"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestCacheHandler_AnswerDNSQuestion(t *testing.T) {
	t.Parallel()
	const (
		aRecordQuestion   = "gitlab.com."
		ptrRecordQuestion = "5.10.0.10.in-addr.arpa"
	)
	type fields struct {
		Cache    *dnsmock.CacheMock
		Fallback dns.Handler
	}
	tests := []struct {
		name     string
		fields   fields
		question dns.Question

		want    any
		wantErr bool
	}{
		{
			name: "Resolve A question - entry from cache",
			fields: fields{
				Cache: &dnsmock.CacheMock{
					OnForwardLookup: func(ctx dnsmock.CacheMockCallsContext, host string) net.IP {
						if host == aRecordQuestion {
							return net.IPv4(10, 0, 10, 5)
						}
						return nil
					},
				},
				Fallback: nil,
			},
			question: dns.Question{
				Name:   "gitlab.com.",
				Qtype:  mdns.TypeA,
				Qclass: mdns.ClassINET,
			},
			want: td.Struct(&mdns.A{
				A: net.IPv4(10, 0, 10, 5),
			}, td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Resolve A question - entry from fallback",
			fields: fields{
				Cache: &dnsmock.CacheMock{
					OnPutRecord: func(ctx dnsmock.CacheMockCallsContext, host string, address net.IP) {
						if len(ctx.PutRecord) > 1 {
							panic(ctx.PutRecord)
						}
					},
				},
				Fallback: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return &mdns.A{
						A: net.IPv4(10, 0, 10, 17),
					}, nil
				}),
			},
			question: dns.Question{
				Name:   aRecordQuestion,
				Qtype:  mdns.TypeA,
				Qclass: mdns.ClassINET,
			},
			want: td.Struct(&mdns.A{
				A: net.IPv4(10, 0, 10, 17),
			}, td.StructFields{}),
			wantErr: false,
		},
		{
			name: "Resolve PTR question - entry in cache",
			fields: fields{
				Cache: &dnsmock.CacheMock{
					OnReverseLookup: func(_ dnsmock.CacheMockCallsContext, _ net.IP) (host string, miss bool) {
						return "gitlab.com.", false
					},
				},
				Fallback: nil,
			},
			question: dns.Question{
				Name:   ptrRecordQuestion,
				Qtype:  mdns.TypePTR,
				Qclass: mdns.ClassINET,
			},
			want: td.Struct(&mdns.PTR{
				Ptr: aRecordQuestion,
			}, td.StructFields{}),
		},
		{
			name: "Don't resolve PTR question - entry not in cache",
			fields: fields{
				Cache: &dnsmock.CacheMock{
					OnReverseLookup: func(_ dnsmock.CacheMockCallsContext, _ net.IP) (host string, miss bool) {
						return "", true
					},
				},
				Fallback: nil,
			},
			question: dns.Question{
				Name:   ptrRecordQuestion,
				Qtype:  mdns.TypePTR,
				Qclass: mdns.ClassINET,
			},
			want:    td.Nil(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const defaultTTL = 5 * time.Second
			h := &dns.CacheHandler{
				Cache:    tt.fields.Cache,
				TTL:      defaultTTL,
				Fallback: tt.fields.Fallback,
			}
			gotRr, err := h.AnswerDNSQuestion(tt.question)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnswerDNSQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotRr, tt.want)
		})
	}
}
