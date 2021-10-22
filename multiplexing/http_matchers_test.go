package multiplexing_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/soheilhy/cmux"

	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/internal/test/integration"
	"gitlab.com/inetmock/inetmock/multiplexing"
)

func TestHTTP(t *testing.T) {
	t.Parallel()
	type args struct {
		behavior []string
		req      *http.Request
	}
	tests := []struct {
		name            string
		args            args
		matcher         cmux.Matcher
		httpClientSetup func(lis test.InMemListener) *http.Client
		wantErr         bool
		want            interface{}
	}{
		{
			name: "Match a HTTP/1.1 request",
			args: args{
				behavior: []string{
					`=> Status(200)`,
				},
				req: &http.Request{
					URL:    test.MustParseURL("https://gitlab.com/inetmock/inetmock"),
					Method: http.MethodGet,
				},
			},
			matcher:         multiplexing.HTTP(),
			httpClientSetup: test.HTTPClientForInMemListener,
			wantErr:         false,
			want: td.Struct(&http.Response{
				StatusCode: 200,
			}, td.StructFields{}),
		},
		{
			name: "Match a HTTP 2 request",
			args: args{
				behavior: []string{
					`=> Status(200)`,
				},
				req: &http.Request{
					URL:    test.MustParseURL("https://gitlab.com/inetmock/inetmock"),
					Method: http.MethodGet,
				},
			},
			matcher:         multiplexing.HTTP(),
			httpClientSetup: test.HTTP2ClientForInMemListener,
			wantErr:         false,
			want: td.Struct(&http.Response{
				StatusCode: 200,
			}, td.StructFields{}),
		},
		{
			name: "Match /dns-query path in request",
			args: args{
				behavior: []string{
					`=> Status(200)`,
				},
				req: &http.Request{
					URL:    test.MustParseURL("https://quad9.com/dns-query"),
					Method: http.MethodGet,
				},
			},
			matcher:         multiplexing.HTTPMatchAnd(multiplexing.PathPrefixMatcher("/dns-query")),
			httpClientSetup: test.HTTPClientForInMemListener,
			wantErr:         false,
			want: td.Struct(&http.Response{
				StatusCode: 200,
			}, td.StructFields{}),
		},
		{
			name: "Don't Match /dnsquery path in request",
			args: args{
				req: &http.Request{
					URL:    test.MustParseURL("https://quad9.com/dnsquery"),
					Method: http.MethodGet,
				},
			},
			matcher:         multiplexing.HTTPMatchAnd(multiplexing.PathPrefixMatcher("/dns-query")),
			httpClientSetup: test.HTTPClientForInMemListener,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inMemListener := test.NewInMemoryListener(t)
			c := cmux.New(inMemListener)
			multiplexerListener := c.Match(tt.matcher)
			go func() {
				if err := c.Serve(); !errors.Is(err, http.ErrServerClosed) {
					t.Logf("Serve() error = %v", err)
				}
			}()
			t.Cleanup(func() {
				c.Close()
			})
			srv := integration.NewTestHTTPServer(t, tt.args.behavior, nil)
			go srv.Listen(t, multiplexerListener)
			client := tt.httpClientSetup(inMemListener)
			resp, err := client.Do(tt.args.req)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("client.Do() error = %v", err)
				}
				return
			}
			td.Cmp(t, resp, tt.want)
		})
	}
}
