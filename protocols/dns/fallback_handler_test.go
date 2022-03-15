package dns_test

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestFallbackHandler_AnswerDNSQuestion(t *testing.T) {
	const defaultTTL = 30 * time.Second
	t.Parallel()
	type fields struct {
		Resolver dns.IPResolver
		Handler  dns.Handler
	}
	tests := []struct {
		name    string
		fields  fields
		want    any
		wantErr bool
	}{
		{
			name: "Get answer from backing handler",
			fields: fields{
				Handler: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return new(mdns.A), nil
				}),
			},
			want: td.Struct(new(mdns.A), td.StructFields{}),
		},
		{
			name: "Handle question with fallback resolver",
			fields: fields{
				Resolver: dns.IPResolverFunc(func(host string) net.IP {
					return net.IPv4(10, 10, 0, 4)
				}),
				Handler: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return nil, dns.ErrNoAnswerForQuestion
				}),
			},
			want: td.Struct(new(mdns.A), td.StructFields{}),
		},
		{
			name: "Neither handler nor fallback can respond to question",
			fields: fields{
				Resolver: dns.IPResolverFunc(func(host string) net.IP {
					return nil
				}),
				Handler: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return nil, dns.ErrNoAnswerForQuestion
				}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := dns.FallbackHandler(tt.fields.Handler, tt.fields.Resolver, defaultTTL)
			got, err := h.AnswerDNSQuestion(dns.Question{Qtype: mdns.TypeA})
			if (err != nil) != tt.wantErr {
				t.Errorf("AnswerDNSQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
