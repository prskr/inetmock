package health

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

//nolint:funlen
func Test_checker_IsHealthy(t *testing.T) {
	t.Parallel()
	type fields struct {
		componentChecks map[string]Check
	}
	type testCase struct {
		name   string
		fields fields
		wantR  Result
	}
	tests := []testCase{
		{
			name: "No checks registered expect HEALTHY",
			fields: fields{
				componentChecks: map[string]Check{},
			},
			wantR: Result{
				Status:     HEALTHY,
				Components: map[string]CheckResult{},
			},
		},
		{
			name: "Return only HEALTHY result expect HEALTHY",
			fields: fields{
				componentChecks: map[string]Check{
					"asdf": func() CheckResult {
						return CheckResult{
							Status:  HEALTHY,
							Message: "",
						}
					},
				},
			},
			wantR: Result{
				Status: HEALTHY,
				Components: map[string]CheckResult{
					"asdf": {
						Status:  HEALTHY,
						Message: "",
					},
				},
			},
		},
		{
			name: "Return HEALTHY and INITIALIZING result expect INITIALIZING",
			fields: fields{
				componentChecks: map[string]Check{
					"asdf": func() CheckResult {
						return CheckResult{
							Status:  HEALTHY,
							Message: "",
						}
					},
					"qwert": func() CheckResult {
						return CheckResult{
							Status:  INITIALIZING,
							Message: "",
						}
					},
				},
			},
			wantR: Result{
				Status: INITIALIZING,
				Components: map[string]CheckResult{
					"asdf": {
						Status:  HEALTHY,
						Message: "",
					},
					"qwert": {
						Status:  INITIALIZING,
						Message: "",
					},
				},
			},
		},
		{
			name: "Return HEALTHY AND UNHEALTHY result expect UNHEALTHY",
			fields: fields{
				componentChecks: map[string]Check{
					"asdf": func() CheckResult {
						return CheckResult{
							Status:  UNHEALTHY,
							Message: "",
						}
					},
					"qwert": func() CheckResult {
						return CheckResult{
							Status:  HEALTHY,
							Message: "",
						}
					},
				},
			},
			wantR: Result{
				Status: UNHEALTHY,
				Components: map[string]CheckResult{
					"asdf": {
						Status:  UNHEALTHY,
						Message: "",
					},
					"qwert": {
						Status:  HEALTHY,
						Message: "",
					},
				},
			},
		},
		{
			name: "Return HEALTHY AND UNKNOWN result expect UNHEALTHY",
			fields: fields{
				componentChecks: map[string]Check{
					"asdf": func() CheckResult {
						return CheckResult{
							Status:  UNKNOWN,
							Message: "",
						}
					},
					"qwert": func() CheckResult {
						return CheckResult{
							Status:  HEALTHY,
							Message: "",
						}
					},
				},
			},
			wantR: Result{
				Status: UNHEALTHY,
				Components: map[string]CheckResult{
					"asdf": {
						Status:  UNKNOWN,
						Message: "",
					},
					"qwert": {
						Status:  HEALTHY,
						Message: "",
					},
				},
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &checker{
				componentChecks: tt.fields.componentChecks,
			}
			gotR := c.IsHealthy()
			td.Cmp(t, gotR, tt.wantR)
		})
	}
}
