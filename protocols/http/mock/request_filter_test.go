package mock_test

import (
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/rules"
	"inetmock.icb4dc0.de/inetmock/protocols/http/mock"
)

func TestHTTPMethodMatcher(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
		req  *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantErr   bool
	}{
		{
			name: "Match GET request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("GET"),
					},
				},
				req: &http.Request{
					Method: http.MethodGet,
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Do not match POST request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("POST"),
					},
				},
				req: &http.Request{
					Method: http.MethodGet,
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name: "Expect error due to argument type mismatch",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
				req: &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				args: []rules.Param{},
				req:  &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.HTTPMethodMatcher(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("HTTPMethodMatcher() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, got.Matches(tt.args.req), tt.wantMatch)
		})
	}
}

func TestPathPatternMatcher(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
		req  *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantErr   bool
	}{
		{
			name: "Match .html request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(".*\\.(?i)htm(l)?$"),
					},
				},
				req: tdhttp.NewRequest(http.MethodGet, "https://www.reddit.com/index.htm", nil),
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Do not match JPGEG request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("POST"),
					},
				},
				req: tdhttp.NewRequest(
					http.MethodGet,
					"https://www.reddit.com/idx.jpeg",
					nil,
				),
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name: "Expect error due to argument type mismatch",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
				req: &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				args: []rules.Param{},
				req:  &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.PathPatternMatcher(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PathPatternMatcher() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, got.Matches(tt.args.req), tt.wantMatch)
		})
	}
}

func TestHeaderValueMatcher(t *testing.T) {
	t.Parallel()
	type args struct {
		args []rules.Param
		req  *http.Request
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantErr   bool
	}{
		{
			name: "Match text/html request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("Accept"),
					},
					{
						String: rules.StringP("text/html"),
					},
				},
				req: &http.Request{
					Header: map[string][]string{
						"Accept": {"text/plain", "text/html"},
					},
				},
			},
			wantMatch: true,
			wantErr:   false,
		},
		{
			name: "Do not match text/plain request",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("Accept"),
					},
					{
						String: rules.StringP("text/html"),
					},
				},
				req: &http.Request{
					Header: map[string][]string{
						"Accept": {"text/plain"},
					},
				},
			},
			wantMatch: false,
			wantErr:   false,
		},
		{
			name: "Expect error due to argument type mismatch",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
					{
						String: rules.StringP("text/html"),
					},
				},
				req: &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
				req: &http.Request{},
			},
			wantMatch: false,
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mock.HeaderValueMatcher(tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("HeaderValueMatcher() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			td.Cmp(t, got.Matches(tt.args.req), tt.wantMatch)
		})
	}
}
