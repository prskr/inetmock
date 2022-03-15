package health_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/pkg/health"
)

func Test_healthHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		checks     []health.Check
		wantStatus int
		wantBody   any
	}{
		{
			name:       "Empty checker - no error",
			wantStatus: 204,
			wantBody:   nil,
		},
		{
			name:       "Checker with success check",
			wantStatus: 204,
			checks: []health.Check{
				health.NewCheckFunc("Success", func(context.Context) error {
					return nil
				}),
			},
			wantBody: nil,
		},
		{
			name:       "Checker with error check",
			wantStatus: 503,
			checks: []health.Check{
				health.NewCheckFunc("Error", func(context.Context) error {
					return errors.New("there's something strange...in the neighborhood")
				}),
			},
			wantBody: td.Map(make(map[string]string), td.MapEntries{
				"Error": td.Contains("something strange"),
			}),
		},
		{
			name:       "Checker with multiple error checks",
			wantStatus: 503,
			checks: []health.Check{
				health.NewCheckFunc("Err1", func(context.Context) error {
					return errors.New("there's something strange...in the neighborhood")
				}),
				health.NewCheckFunc("Err2", func(context.Context) error {
					return errors.New("who you gonna call")
				}),
			},
			wantBody: td.Map(make(map[string]string), td.MapEntries{
				"Err1": td.Contains("something strange"),
				"Err2": td.Contains("gonna call"),
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checker := health.New()

			for idx := range tt.checks {
				check := tt.checks[idx]
				if err := checker.AddCheck(check); err != nil {
					t.Errorf("checker.AddCheck() error = %v", err)
					return
				}
			}

			h := health.NewHealthHandler(checker)
			ta := tdhttp.NewTestAPI(t, h)
			ta = ta.Get("/").
				CmpStatus(tt.wantStatus)

			if tt.wantBody != nil {
				ta.
					CmpHeader(td.SuperMapOf(http.Header{}, td.MapEntries{
						"Content-Type": []string{"application/json"},
					})).
					CmpJSONBody(tt.wantBody)
			}
		})
	}
}
