package dns_test

import (
	"net"
	"testing"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/dns"
)

func TestResponseFilters(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
		resp *dns.Response
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		wantMatchErr      bool
		validatorProvider func(args ...rules.Param) (dns.Validator, error)
	}{
		{
			name: "NotEmpty - nil result expect matcher to error",
			args: args{
				resp: nil,
			},
			validatorProvider: dns.NotEmtpyResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "NotEmpty - Empty result expect matcher to error",
			args: args{
				resp: new(dns.Response),
			},
			validatorProvider: dns.NotEmtpyResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "NotEmpty - Hosts set - no error",
			args: args{
				resp: &dns.Response{
					Hosts: []string{"google.com"},
				},
			},
			validatorProvider: dns.NotEmtpyResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "NotEmpty - Addresses set - no error",
			args: args{
				resp: &dns.Response{
					Addresses: []net.IP{net.IPv4(192, 168, 0, 1)},
				},
			},
			validatorProvider: dns.NotEmtpyResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedHost - Missing host name param",
			args: args{
				args: make([]rules.Param, 0),
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedHost - nil param",
			args: args{
				args: make([]rules.Param, 1),
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedHost - Wrong param type",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedHost - Response nil - expect error",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resp: nil,
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedHost - Response completely empty - expect error",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resp: new(dns.Response),
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedHost - Response hosts empty",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resp: &dns.Response{
					Addresses: make([]net.IP, 1),
				},
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedHost - Response does not match",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resp: &dns.Response{
					Hosts: []string{"about.gitlab.com."},
				},
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedHost - First response matches",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("gitlab.com"),
					},
				},
				resp: &dns.Response{
					Hosts: []string{
						"gitlab.com.",
					},
				},
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedHost - Second response matches",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("about.gitlab.com"),
					},
				},
				resp: &dns.Response{
					Hosts: []string{
						"gitlab.com.",
						"about.gitlab.com.",
					},
				},
			},
			validatorProvider: dns.ResolvedHostResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},

		{
			name: "ResolvedIP - Missing host name param",
			args: args{
				args: make([]rules.Param, 0),
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedIP - nil param",
			args: args{
				args: make([]rules.Param, 1),
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedIP - Wrong param type",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedIP - Response nil - expect error",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: nil,
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedIP - Response completely empty - expect error",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: new(dns.Response),
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedIP - Response addresses empty",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: &dns.Response{
					Hosts: make([]string, 1),
				},
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedIP - Response does not match",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(192, 168, 1, 1),
					},
				},
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "ResolvedIP - First response matches",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(192, 168, 0, 11),
					},
				},
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "ResolvedIP - Second response matches",
			args: args{
				args: []rules.Param{
					{
						IP: net.IPv4(192, 168, 0, 11),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(192, 168, 0, 42),
						net.IPv4(192, 168, 0, 11),
					},
				},
			},
			validatorProvider: dns.ResolvedIPResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "InCIDR - Missing CIDR param",
			args: args{
				args: make([]rules.Param, 0),
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "InCIDR - nil param",
			args: args{
				args: make([]rules.Param, 1),
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "InCIDR - Wrong param type",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           true,
			wantMatchErr:      false,
		},
		{
			name: "InCIDR - Response nil - expect error",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: nil,
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "InCIDR - Response completely empty - expect error",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: new(dns.Response),
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "InCIDR - Response addresses empty",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: &dns.Response{
					Hosts: make([]string, 1),
				},
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "InCIDR - Response does not match",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(192, 168, 1, 1),
					},
				},
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      true,
		},
		{
			name: "InCIDR - First response matches",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(10, 1, 0, 10),
					},
				},
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
		{
			name: "InCIDR - Second response matches",
			args: args{
				args: []rules.Param{
					{
						CIDR: rules.MustParseCIDR("10.1.0.0/16"),
					},
				},
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(192, 168, 0, 42),
						net.IPv4(10, 1, 0, 10),
					},
				},
			},
			validatorProvider: dns.InCIDRResponseFilter,
			wantErr:           false,
			wantMatchErr:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var validator, err = tt.validatorProvider(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvedHostResponseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if validator == nil {
				if !tt.wantErr {
					t.Error("returned validator is nil")
				}
				return
			}

			if err := validator.Matches(tt.args.resp); (err != nil) != tt.wantMatchErr {
				t.Errorf("validator.Matches() error = %v, wantMatchErr = %t", err, tt.wantMatchErr)
			}
		})
	}
}

func TestValidatorsForRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule string
		resp *dns.Response
	}
	tests := []struct {
		name            string
		args            args
		wantChainLength int
		wantErr         bool
		wantMatchErr    bool
	}{
		{
			name: "Empty filter chain",
			args: args{
				rule: `A("gitlab.com")`,
			},
		},
		{
			name: "Rule parsing error",
			args: args{
				rule: `A("gitlab.com)`,
			},
			wantErr: true,
		},
		{
			name: "unmatched filter",
			args: args{
				rule: `A("gitlab.com") => Fuck()`,
			},
			wantErr: true,
		},
		{
			name: "Parse NotEmpty filter and match without error",
			args: args{
				rule: `A("gitlab.com") => NotEmpty()`,
				resp: &dns.Response{
					Addresses: make([]net.IP, 1),
				},
			},
			wantChainLength: 1,
			wantErr:         false,
			wantMatchErr:    false,
		},
		{
			name: "Parse ResolvedIP filter and match without error",
			args: args{
				rule: `A("gitlab.com") => ResolvedIP(1.1.1.1)`,
				resp: &dns.Response{
					Addresses: []net.IP{
						net.IPv4(1, 1, 1, 1),
					},
				},
			},
			wantChainLength: 1,
			wantErr:         false,
			wantMatchErr:    false,
		},
		{
			name: "Parse ResolvedHost filter and match without error",
			args: args{
				rule: `PTR(1.1.1.1) => ResolvedHost("one.one.one.one")`,
				resp: &dns.Response{
					Hosts: []string{"one.one.one.one."},
				},
			},
			wantChainLength: 1,
			wantErr:         false,
			wantMatchErr:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var rule = new(rules.Check)
			if err := rules.Parse(tt.args.rule, rule); err != nil {
				if !tt.wantErr {
					t.Errorf("rules.Parse() error = %v", err)
				}
				return
			}

			chain, err := dns.ValidatorsForRule(rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatorsForRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantChainLength != chain.Len() {
				t.Errorf("Chain has length %d but want length %d", len(chain), tt.wantChainLength)
				return
			}

			if err := chain.Matches(tt.args.resp); (err != nil) != tt.wantMatchErr {
				t.Errorf("chain.Matches() error = %v, wantMatchErr = %t", err, tt.wantMatchErr)
			}
		})
	}
}

func TestValidationChain_Matches(t *testing.T) {
	t.Parallel()
	type args struct {
		resp *dns.Response
	}
	tests := []struct {
		name       string
		chainSetup func(tb testing.TB) dns.ValidationChain
		args       args
		wantErr    bool
	}{
		{
			name: "nil chain",
			chainSetup: func(tb testing.TB) dns.ValidationChain {
				tb.Helper()
				return nil
			},
			wantErr: false,
		},
		{
			name: "Empty chain",
			chainSetup: func(tb testing.TB) dns.ValidationChain {
				tb.Helper()
				return make(dns.ValidationChain, 0)
			},
			wantErr: false,
		},
		{
			name: "Matching chain",
			chainSetup: func(tb testing.TB) dns.ValidationChain {
				tb.Helper()
				if validator, err := dns.NotEmtpyResponseFilter(); err != nil {
					tb.Errorf("dns.NotEmtpyResponseFilter() error = %v", err)
					return nil
				} else {
					return dns.ValidationChain{
						validator,
					}
				}
			},
			args: args{
				resp: &dns.Response{
					Hosts: make([]string, 1),
				},
			},
		},
		{
			name: "Not matching chain",
			chainSetup: func(tb testing.TB) dns.ValidationChain {
				tb.Helper()
				if validator, err := dns.NotEmtpyResponseFilter(); err != nil {
					tb.Errorf("dns.NotEmtpyResponseFilter() error = %v", err)
					return nil
				} else {
					return dns.ValidationChain{
						validator,
					}
				}
			},
			args: args{
				resp: &dns.Response{
					Hosts: make([]string, 0),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var chain = tt.chainSetup(t)
			if err := chain.Matches(tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("ValidationChain.Matches() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationChain_Add(t *testing.T) {
	t.Parallel()
	type args struct {
		v dns.Validator
	}
	tests := []struct {
		name            string
		chain           dns.ValidationChain
		args            args
		wantFinalLength int
	}{
		{
			name:  "Empty chain",
			chain: nil,
			args: args{
				v: nil,
			},
			wantFinalLength: 1,
		},
		{
			name: "Non-empty chain",
			args: args{
				v: nil,
			},
			chain: func() dns.ValidationChain {
				f, err := dns.NotEmtpyResponseFilter()
				if err != nil {
					panic(err)
				}
				return dns.ValidationChain{f}
			}(),
			wantFinalLength: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.chain.Add(tt.args.v)
			if newLength := tt.chain.Len(); tt.wantFinalLength != newLength {
				t.Errorf("Current length %d did not match actual length %d", newLength, tt.wantFinalLength)
			}
		})
	}
}
