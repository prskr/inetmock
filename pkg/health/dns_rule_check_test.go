package health_test

import (
	"context"
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health"
	"gitlab.com/inetmock/inetmock/pkg/health/dns"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func TestNewDNSRuleCheck(t *testing.T) {
	t.Parallel()
	type args struct {
		name     string
		resolver health.ResolverForModule
		check    string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantParseErr bool
		wantCheckErr bool
	}{
		{
			name: "Empty name",
			args: args{
				name: "",
			},
			wantErr:      true,
			wantParseErr: true,
		},
		{
			name: "Resolver nil",
			args: args{
				name:     "test",
				resolver: nil,
			},
			wantErr:      true,
			wantParseErr: true,
		},
		{
			name: "No error get non-nil value check",
			args: args{
				name:     "test",
				resolver: new(dns.MockResolver),
				check:    `dns.A("gitlab.com")`,
			},
			wantErr: false,
		},
		{
			name: "Mocked resolver expect match error",
			args: args{
				name: "test",
				resolver: &dns.MockResolver{
					LookupHostDelegate: func(context.Context, string) (addrs []net.IP, err error) {
						return []net.IP{net.IPv4(192, 168, 0, 1)}, nil
					},
				},
				check: `dns.A("api.inetmock.loc") => ResolvedIP(192.168.1.42)`,
			},
			wantErr:      false,
			wantParseErr: false,
			wantCheckErr: true,
		},
		{
			name: "Mocked resolver expect no match error",
			args: args{
				name: "test",
				resolver: &dns.MockResolver{
					LookupHostDelegate: func(context.Context, string) (addrs []net.IP, err error) {
						return []net.IP{net.IPv4(192, 168, 1, 42)}, nil
					},
				},
				check: `dns.A("api.inetmock.loc") => ResolvedIP(192.168.1.42)`,
			},
			wantErr:      false,
			wantParseErr: false,
			wantCheckErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parsedCheck := new(rules.Check)
			if err := rules.Parse(tt.args.check, parsedCheck); (err != nil) != tt.wantParseErr {
				t.Errorf("rules.Parse() error = %v", err)
				return
			}
			logger := logging.CreateTestLogger(t)
			compiledCheck, err := health.NewDNSRuleCheck(tt.args.name, tt.args.resolver, logger, parsedCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDNSRuleCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if compiledCheck == nil {
				if !tt.wantErr {
					t.Error("compiled check is nil but no error was expected")
				}
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			if err := compiledCheck.Status(ctx); (err != nil) != tt.wantCheckErr {
				t.Errorf("compiledCheck.Status() error = %v, wantCheckErr = %t", err, tt.wantCheckErr)
			}
		})
	}
}
