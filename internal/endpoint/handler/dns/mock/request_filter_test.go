package mock_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/miekg/dns"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	"gitlab.com/inetmock/inetmock/internal/rules"
)

func TestHostnameQuestionFilter(t *testing.T) {
	t.Parallel()
	type args struct {
		qType  uint16
		req    *dns.Question
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
				qType: dns.TypeA,
				req: &dns.Question{
					Qtype: dns.TypeA,
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
				qType: dns.TypeA,
				req: &dns.Question{
					Qtype: dns.TypeA,
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
				qType: dns.TypeAAAA,
				req: &dns.Question{
					Qtype: dns.TypeAAAA,
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
				qType: dns.TypeAAAA,
				req: &dns.Question{
					Qtype: dns.TypeSRV,
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
				qType: dns.TypeAAAA,
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "Fail to compile pattern",
			args: args{
				qType: dns.TypeAAAA,
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
			got, err := mock.HostnameQuestionFilter(tt.args.qType)(tt.args.params...)
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

func TestRequestFiltersForRoutingRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule *rules.Routing
	}
	tests := []struct {
		name        string
		args        args
		wantFilters interface{}
		wantErr     bool
	}{
		{
			name: "Unknown filter method",
			args: args{
				rule: &rules.Routing{
					Filters: &rules.Filters{
						Chain: []rules.Method{
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
			name: "Filters empty",
			args: args{
				rule: &rules.Routing{
					Filters: &rules.Filters{
						Chain: make([]rules.Method, 0),
					},
				},
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "Filters nil",
			args: args{
				rule: new(rules.Routing),
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "Routing nil",
			args: args{
				rule: new(rules.Routing),
			},
			wantFilters: td.Nil(),
			wantErr:     false,
		},
		{
			name: "Single A filter",
			args: args{
				rule: &rules.Routing{
					Filters: &rules.Filters{
						Chain: []rules.Method{
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
			wantFilters: td.Code(func(filters []mock.RequestFilter) bool {
				if len(filters) != 1 {
					return false
				}
				return filters[0].Matches(&dns.Question{
					Qtype: dns.TypeA,
					Name:  "www.google.com",
				})
			}),
			wantErr: false,
		},
		{
			name: "Single AAAA filter",
			args: args{
				rule: &rules.Routing{
					Filters: &rules.Filters{
						Chain: []rules.Method{
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
			wantFilters: td.Code(func(filters []mock.RequestFilter) bool {
				if len(filters) != 1 {
					return false
				}
				return filters[0].Matches(&dns.Question{
					Qtype: dns.TypeAAAA,
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
			gotFilters, err := mock.RequestFiltersForRoutingRule(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestFiltersForRoutingRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, gotFilters, tt.wantFilters)
		})
	}
}
