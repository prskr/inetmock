package dns_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/protocols/dns"
)

func TestHostnameQuestionFilter(t *testing.T) {
	t.Parallel()
	type args struct {
		qType  uint16
		req    dns.Question
		params []rules.Param
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantErr   bool
	}{
		{
			name: "Match A request",
			args: args{
				qType: mdns.TypeA,
				req: dns.Question{
					Qtype: mdns.TypeA,
					Name:  "google.com",
				},
				params: []rules.Param{
					{
						String: rules.StringP(".*"),
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Match A request with pattern",
			args: args{
				qType: mdns.TypeA,
				req: dns.Question{
					Qtype: mdns.TypeA,
					Name:  `.*google\.com`,
				},
				params: []rules.Param{
					{
						String: rules.StringP(".*"),
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Match AAAA request",
			args: args{
				qType: mdns.TypeAAAA,
				req: dns.Question{
					Qtype: mdns.TypeAAAA,
					Name:  "google.com",
				},
				params: []rules.Param{
					{
						String: rules.StringP(".*"),
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Don't match SRV request",
			args: args{
				qType: mdns.TypeAAAA,
				req: dns.Question{
					Qtype: mdns.TypeSRV,
					Name:  "google.com",
				},
				params: []rules.Param{
					{
						String: rules.StringP(".*"),
					},
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name: "Fail due to missing parameter",
			args: args{
				qType: mdns.TypeAAAA,
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "Fail to compile pattern",
			args: args{
				qType: mdns.TypeAAAA,
				params: []rules.Param{
					{
						String: rules.StringP(`[`),
					},
				},
			},
			wantMatch: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := dns.HostnameQuestionFilter(tt.args.qType)(tt.args.params...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("HostnameQuestionFilter() err = %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Expected to fail but did not")
			}

			if got.Matches(tt.args.req) != tt.wantMatch {
				t.Errorf("Expected to match request but did not")
			}
		})
	}
}

func TestQuestionPredicatesForRoutingRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule *rules.SingleResponsePipeline
	}
	tests := []struct {
		name        string
		args        args
		wantFilters any
		wantErr     bool
	}{
		{
			name: "Unknown filter method",
			args: args{
				rule: &rules.SingleResponsePipeline{
					FilterChain: &rules.Filters{
						Chain: []rules.Call{
							{
								Name: "srv",
							},
						},
					},
				},
			},
			wantFilters: td.Empty(),
			wantErr:     true,
		},
		{
			name: "FilterChain empty",
			args: args{
				rule: &rules.SingleResponsePipeline{
					FilterChain: &rules.Filters{
						Chain: make([]rules.Call, 0),
					},
				},
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "FilterChain nil",
			args: args{
				rule: new(rules.SingleResponsePipeline),
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "SingleResponsePipeline nil",
			args: args{
				rule: new(rules.SingleResponsePipeline),
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "Single A filter",
			args: args{
				rule: &rules.SingleResponsePipeline{
					FilterChain: &rules.Filters{
						Chain: []rules.Call{
							{
								Name: "A",
								Params: []rules.Param{
									{
										String: rules.StringP(`.*google\.com`),
									},
								},
							},
						},
					},
				},
			},
			wantFilters: td.Code(func(filters []dns.QuestionPredicate) bool {
				if len(filters) != 1 {
					return false
				}
				return filters[0].Matches(dns.Question{
					Qtype: mdns.TypeA,
					Name:  "www.google.com",
				})
			}),
			wantErr: false,
		},
		{
			name: "Single AAAA filter",
			args: args{
				rule: &rules.SingleResponsePipeline{
					FilterChain: &rules.Filters{
						Chain: []rules.Call{
							{
								Name: "AAAA",
								Params: []rules.Param{
									{
										String: rules.StringP(`.*google\.com`),
									},
								},
							},
						},
					},
				},
			},
			wantFilters: td.Code(func(filters []dns.QuestionPredicate) bool {
				if len(filters) != 1 {
					return false
				}
				return filters[0].Matches(dns.Question{
					Qtype: mdns.TypeAAAA,
					Name:  "www.google.com",
				})
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotFilters, err := dns.QuestionPredicatesForRoutingRule(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestFiltersForRoutingRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotFilters, tt.wantFilters)
		})
	}
}
