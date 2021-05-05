package health_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/health"
)

func Test_healthHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		checks   []health.Check
		want     interface{}
		wantBody interface{}
	}{
		{
			name: "Empty checker - no error",
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": td.Between(200, 299),
			}),
			wantBody: "",
		},
		{
			name: "Checker with success check",
			checks: []health.Check{
				health.NewCheckFunc("Success", func(context.Context) error {
					return nil
				}),
			},
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": td.Between(200, 299),
			}),
			wantBody: "",
		},
		{
			name: "Checker with error check",
			checks: []health.Check{
				health.NewCheckFunc("Error", func(context.Context) error {
					return errors.New("there's something strange...in the neighborhood")
				}),
			},
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": 503,
				"Header": td.SuperMapOf(http.Header{}, td.MapEntries{
					"Content-Type": td.Contains("application/json"),
				}),
			}),
			wantBody: td.Contains("something strange"),
		},
		{
			name: "Checker with multiple error checks",
			checks: []health.Check{
				health.NewCheckFunc("Err1", func(context.Context) error {
					return errors.New("there's something strange...in the neighborhood")
				}),
				health.NewCheckFunc("Err2", func(context.Context) error {
					return errors.New("who you gonna call")
				}),
			},
			want: td.Struct(new(http.Response), td.StructFields{
				"StatusCode": 503,
				"Header": td.SuperMapOf(http.Header{}, td.MapEntries{
					"Content-Type": td.Contains("application/json"),
				}),
			}),
			wantBody: td.Contains("something strange"),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var checker = health.New()

			for idx := range tt.checks {
				check := tt.checks[idx]
				if err := checker.AddCheck(check); err != nil {
					t.Errorf("checker.AddCheck() error = %v", err)
					return
				}
			}

			var h = health.NewHealthHandler(checker)
			var listener = eptest.NewInMemoryListener(t)

			go func() {
				if err := http.Serve(listener, h); err != nil && !errors.Is(err, http.ErrServerClosed) {
					t.Errorf("http.Serve() error = %v", err)
				}
			}()

			var client = eptest.HTTPClientForInMemListener(listener)
			ctx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/", nil)
			if resp, err := client.Do(req); err != nil {
				t.Errorf("failed to fetch health state = error %v", err)
				return
			} else {
				var bodyBuilder = new(strings.Builder)
				_, _ = io.Copy(bodyBuilder, resp.Body)
				defer func() {
					_ = resp.Body.Close()
				}()
				td.Cmp(t, resp, tt.want)
				td.Cmp(t, bodyBuilder.String(), tt.wantBody)
			}
		})
	}
}
