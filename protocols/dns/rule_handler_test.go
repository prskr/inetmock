package dns_test

import (
	"net"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestRuleHandler_AnswerDNSQuestion(t *testing.T) {
	t.Parallel()
	const defaultTTL = 30 * time.Second
	tests := []struct {
		name     string
		rawRules []string
		question dns.Question
		want     interface{}
		wantErr  bool
	}{
		{
			name:    "No rule expect error",
			wantErr: true,
		},
		{
			name: "Rule without filter",
			rawRules: []string{
				`=> IP(1.1.1.1)`,
			},
			question: dns.Question{Qtype: mdns.TypeA},
			want: td.Struct(&mdns.A{
				A: net.IPv4(1, 1, 1, 1),
			}, td.StructFields{}),
		},
		{
			name: "Rule with matching filter",
			rawRules: []string{
				`A("gitlab.com") => IP(1.2.3.4)`,
			},
			question: dns.Question{Qtype: mdns.TypeA, Name: "gitlab.com"},
			want: td.Struct(&mdns.A{
				A: net.IPv4(1, 2, 3, 4),
			}, td.StructFields{}),
		},
		{
			name: "Rule with not matching filter",
			rawRules: []string{
				`A("gitlab.com") => IP(1.2.3.4)`,
			},
			question: dns.Question{Qtype: mdns.TypeA, Name: "github.com"},
			wantErr:  true,
		},
		{
			name: "Rule with not matching question type",
			rawRules: []string{
				`A("gitlab.com") => IP(1.2.3.4)`,
			},
			question: dns.Question{Qtype: mdns.TypeAAAA, Name: "gitlab.com"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := dns.RuleHandler{TTL: defaultTTL}
			for _, rawRule := range tt.rawRules {
				if err := r.RegisterRule(rawRule); err != nil {
					t.Errorf("Failed to register rule: %v", err)
					return
				}
			}
			got, err := r.AnswerDNSQuestion(tt.question)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnswerDNSQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
