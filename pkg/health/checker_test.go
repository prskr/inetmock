package health

import (
	"reflect"
	"testing"
)

//nolint:funlen
func Test_checker_IsHealthy(t *testing.T) {
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			c := &checker{
				componentChecks: tt.fields.componentChecks,
			}
			if gotR := c.IsHealthy(); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("IsHealthy() = %v, want %v", gotR, tt.wantR)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
