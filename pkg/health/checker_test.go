package health_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/health"
)

func Test_checker_AddCheck(t *testing.T) {
	t.Parallel()
	type args struct {
		check health.Check
	}
	tests := []struct {
		name         string
		checkerSetup func(tb testing.TB, checker health.Checker)
		args         args
		wantErr      bool
	}{
		{
			name: "Adding check to empty checker",
			args: args{
				check: health.NewCheckFunc("Empty", nil),
			},
			wantErr: false,
		},
		{
			name: "Adding check to non-empty checker",
			checkerSetup: func(tb testing.TB, checker health.Checker) {
				tb.Helper()
				if err := checker.AddCheck(health.NewCheckFunc("Redis", nil)); !td.CmpNoError(tb, err) {
					tb.Fail()
				}
			},
			args: args{
				check: health.NewCheckFunc("MySQL", nil),
			},
			wantErr: false,
		},
		{
			name: "Adding conflicting check",
			checkerSetup: func(tb testing.TB, checker health.Checker) {
				tb.Helper()
				if err := checker.AddCheck(health.NewCheckFunc("Redis", nil)); !td.CmpNoError(tb, err) {
					tb.Fail()
				}
			},
			args: args{
				check: health.NewCheckFunc("Redis", nil),
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checker := health.New()
			if tt.checkerSetup != nil {
				tt.checkerSetup(t, checker)
			}
			if err := checker.AddCheck(tt.args.check); (err != nil) != tt.wantErr {
				t.Errorf("AddCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checker_Status(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		checker health.Checker
		wantRes interface{}
	}{
		{
			name:    "Get status of empty checker - expect empty result",
			checker: health.New(),
			wantRes: health.Result{},
		},
		{
			name: "Get status of single check",
			checker: func() health.Checker {
				checker := health.New()
				_ = checker.AddCheck(newCheckOfResult("Redis", nil))
				return checker
			}(),
			wantRes: health.Result{
				"Redis": nil,
			},
		},
		{
			name: "Get status of multiple checks",
			checker: func() health.Checker {
				checker := health.New()
				_ = checker.AddCheck(newCheckOfResult("Redis", nil))
				_ = checker.AddCheck(newCheckOfResult("MySQL", nil))
				return checker
			}(),
			wantRes: td.Map(health.Result{}, map[interface{}]interface{}{
				"MySQL": nil,
				"Redis": nil,
			}),
		},
		{
			name: "Get status of multiple checks with one error",
			checker: func() health.Checker {
				checker := health.New()
				_ = checker.AddCheck(newCheckOfResult("Redis", nil))
				_ = checker.AddCheck(newCheckOfResult("MySQL", nil))
				_ = checker.AddCheck(newCheckOfResult("HTTP", errors.New("there's something strange in the neighborhood")))
				return checker
			}(),
			wantRes: td.Map(health.Result{}, map[interface{}]interface{}{
				"MySQL": nil,
				"Redis": nil,
				"HTTP":  errors.New("there's something strange in the neighborhood"),
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testCtx, cancel := context.WithTimeout(test.Context(t), 50*time.Millisecond)
			t.Cleanup(cancel)
			gotRes := tt.checker.Status(testCtx)
			td.Cmp(t, gotRes, tt.wantRes)
		})
	}
}

type checkOfResult struct {
	name   string
	result error
}

func newCheckOfResult(name string, result error) health.Check {
	return &checkOfResult{
		name:   name,
		result: result,
	}
}

func (c checkOfResult) Name() string {
	return c.name
}

func (c checkOfResult) Status(context.Context) error {
	return c.result
}
