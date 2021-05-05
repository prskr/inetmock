package http_test

import (
	"io"
	gohttp "net/http"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/health/http"
)

const (
	//nolint:lll
	loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Non tellus orci ac auctor. Mi eget mauris pharetra et ultrices neque ornare aenean. Vitae proin sagittis nisl rhoncus mattis rhoncus urna. Malesuada fames ac turpis egestas sed tempus. Cras ornare arcu dui vivamus arcu. Et tortor consequat id porta nibh venenatis cras sed. Porttitor leo a diam sollicitudin tempor id. Volutpat sed cras ornare arcu dui. Facilisis magna etiam tempor orci. Morbi tincidunt ornare massa eget egestas. Varius sit amet mattis vulputate enim nulla. Aenean et tortor at risus viverra adipiscing at in tellus. Consectetur adipiscing elit ut aliquam purus. Facilisis magna etiam tempor orci eu. Vitae purus faucibus ornare suspendisse sed nisi.`
)

func TestCheckFilter(t *testing.T) {
	t.Parallel()
	type args struct {
		args        []rules.Param
		respToMatch *gohttp.Response
	}
	tests := []struct {
		name           string
		filterProvider func(args ...rules.Param) (http.Validator, error)
		args           args
		wantMatch      bool
		wantErr        bool
	}{
		{
			name:           "StatusCodeFilter - Matching code 200",
			filterProvider: http.StatusCodeFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(200),
					},
				},
				respToMatch: &gohttp.Response{
					StatusCode: 200,
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "StatusCodeFilter - Matching code 204",
			filterProvider: http.StatusCodeFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(204),
					},
				},
				respToMatch: &gohttp.Response{
					StatusCode: 204,
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "StatusCodeFilter - Not matching code 301",
			filterProvider: http.StatusCodeFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(204),
					},
				},
				respToMatch: &gohttp.Response{
					StatusCode: 301,
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "StatusCodeFilter - Not matching nil response",
			filterProvider: http.StatusCodeFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(204),
					},
				},
				respToMatch: nil,
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "StatusCodeFilter - Error due to param mismatch",
			filterProvider: http.StatusCodeFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("204"),
					},
				},
			},
			wantErr: true,
		},
		{
			name:           "StatusCodeFilter - Error due to missing parameter",
			filterProvider: http.StatusCodeFilter,
			wantErr:        true,
		},
		{
			name:           "ResponseBodyContainsFilter - Find a normal string",
			filterProvider: http.ResponseBodyContainsFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("vivamus"),
					},
				},
				respToMatch: &gohttp.Response{
					Body: io.NopCloser(strings.NewReader(loremIpsum)),
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyContainsFilter - Find a JSON substring",
			filterProvider: http.ResponseBodyContainsFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`"lastName":"Tester"`),
					},
				},
				respToMatch: &gohttp.Response{
					Body: io.NopCloser(strings.NewReader(`{"firstName":"Ted","lastName":"Tester"}`)),
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyContainsFilter - Not matching a non-contained string",
			filterProvider: http.ResponseBodyContainsFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`LastName`),
					},
				},
				respToMatch: &gohttp.Response{
					Body: io.NopCloser(strings.NewReader(`{"firstName":"Ted","lastName":"Tester"}`)),
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyContainsFilter - Not matching nil response",
			filterProvider: http.ResponseBodyContainsFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("hello, world"),
					},
				},
				respToMatch: nil,
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyContainsFilter - Error due to param mismatch",
			filterProvider: http.ResponseBodyContainsFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(301),
					},
				},
			},
			wantErr: true,
		},
		{
			name:           "ResponseBodyContainsFilter - Error due to missing parameter",
			filterProvider: http.ResponseBodyContainsFilter,
			wantErr:        true,
		},
		{
			name:           "ResponseHeaderFilter - Match JSON content-type",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("Content-Type"),
					},
					{
						String: rules.StringP("application/json"),
					},
				},
				respToMatch: &gohttp.Response{
					Header: gohttp.Header{
						"Content-Type": []string{"application/json"},
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseHeaderFilter - Match JSON content-type - header name case insensitive",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("content-type"),
					},
					{
						String: rules.StringP("application/json"),
					},
				},
				respToMatch: &gohttp.Response{
					Header: gohttp.Header{
						"Content-Type": []string{"application/json"},
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseHeaderFilter - Match JSON content-type - header value case insensitive",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("content-type"),
					},
					{
						String: rules.StringP("Application/JSON"),
					},
				},
				respToMatch: &gohttp.Response{
					Header: gohttp.Header{
						"Content-Type": []string{"aPplicAtion/json"},
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseHeaderFilter - Match JSON content-type - header value contains search value",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("content-type"),
					},
					{
						String: rules.StringP("Application/JSON"),
					},
				},
				respToMatch: &gohttp.Response{
					Header: gohttp.Header{
						"Content-Type": []string{"aPplicAtion/json; encoding=utf8"},
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseHeaderFilter - Not matching nil response",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("Content-Type"),
					},
					{
						String: rules.StringP("application/json"),
					},
				},
				respToMatch: nil,
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "ResponseHeaderFilter - Error due to param mismatch",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(301),
					},
				},
			},
			wantErr: true,
		},
		{
			name:           "ResponseHeaderFilter - Error due to missing parameter",
			filterProvider: http.ResponseHeaderFilter,
			wantErr:        true,
		},
		{
			name:           "ResponseHeaderFilter - Error due to parameter type mismatch",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
					{
						String: rules.StringP("application/json"),
					},
				},
				respToMatch: nil,
			},
			wantErr: true,
		},
		{
			name:           "ResponseHeaderFilter - Error due to parameter type mismatch",
			filterProvider: http.ResponseHeaderFilter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("application/json"),
					},
					{
						Int: rules.IntP(42),
					},
				},
				respToMatch: nil,
			},
			wantErr: true,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - Matching payload",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("c39925a57c5bf7d91a7f1a1001d58a9aed7f9d158e9638925c175afc11288215"),
					},
				},
				respToMatch: &gohttp.Response{
					Body: io.NopCloser(strings.NewReader(loremIpsum)),
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - not matching payload",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("d39925a57c5bf7d91a7f1a1001d58a9aed7f9d158e9638925c175afc11288215"),
					},
				},
				respToMatch: &gohttp.Response{
					Body: io.NopCloser(strings.NewReader(loremIpsum)),
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - Not matching nil response",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("d39925a57c5bf7d91a7f1a1001d58a9aed7f9d158e9638925c175afc11288215"),
					},
				},
				respToMatch: nil,
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - Error due to invalid hex string",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("hello world"),
					},
				},
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - Error due to missing parameter",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			wantErr:        true,
		},
		{
			name:           "ResponseBodyHashSHA256Filter - Error due to parameter type mismatch",
			filterProvider: http.ResponseBodyHashSHA256Filter,
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
				respToMatch: nil,
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.filterProvider(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StatusCodeFilter() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if matchErr := got.Matches(tt.args.respToMatch); (matchErr == nil) != tt.wantMatch {
				t.Errorf("got.Matches() = %v", matchErr)
			}
		})
	}
}

func TestCheckFiltersForRule(t *testing.T) {
	t.Parallel()
	type args struct {
		rule *rules.Check
	}
	tests := []struct {
		name        string
		args        args
		wantFilters interface{}
		wantErr     bool
	}{
		{
			name: "Empty array if Validators nil",
			args: args{
				rule: new(rules.Check),
			},
			wantFilters: td.Empty(),
			wantErr:     false,
		},
		{
			name: "Empty array for no filters",
			args: args{
				rule: &rules.Check{
					Validators: new(rules.Filters),
				},
			},
			wantFilters: td.Empty(),
			wantErr:     false,
		},
		{
			name: "Not empty array for valid validator name",
			args: args{
				rule: &rules.Check{
					Validators: &rules.Filters{
						Chain: []rules.Method{
							{
								Module: "http",
								Name:   "status",
								Params: []rules.Param{
									{
										Int: rules.IntP(200),
									},
								},
							},
						},
					},
				},
			},
			wantFilters: td.NotEmpty(),
			wantErr:     false,
		},
		{
			name: "Error due to param mismatch",
			args: args{
				rule: &rules.Check{
					Validators: &rules.Filters{
						Chain: []rules.Method{
							{
								Module: "http",
								Name:   "statuscode",
								Params: []rules.Param{
									{
										String: rules.StringP("200"),
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error due to unknown validator",
			args: args{
				rule: &rules.Check{
					Validators: &rules.Filters{
						Chain: []rules.Method{
							{
								Module: "http",
								Name:   "statcode",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotFilters, err := http.ValidatorsForRule(tt.args.rule)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ValidatorsForRule() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, gotFilters, tt.wantFilters)
		})
	}
}
