package health_test

import (
	"errors"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/pkg/health"
)

func TestResult_IsHealthy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		result      health.Result
		wantHealthy bool
	}{
		{
			name:        "Empty expect - expect healthy",
			result:      health.Result{},
			wantHealthy: true,
		},
		{
			name: "Successful test - expect healthy",
			result: health.Result{
				"Sample check": nil,
			},
			wantHealthy: true,
		},
		{
			name: "Failed test - expect unhealthy",
			result: health.Result{
				"Failed check": errors.New("any kind of error"),
			},
			wantHealthy: false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotHealthy := tt.result.IsHealthy(); gotHealthy != tt.wantHealthy {
				t.Errorf("IsHealthy() = %v, want %v", gotHealthy, tt.wantHealthy)
			}
		})
	}
}

func TestResult_CheckResult(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name           string
		result         health.Result
		args           args
		wantKnownCheck bool
		wantErr        bool
	}{
		{
			name: "Known, successful check",
			result: health.Result{
				"Redis": nil,
			},
			args: args{
				"Redis",
			},
			wantKnownCheck: true,
			wantErr:        false,
		},
		{
			name: "Known, failed check",
			result: health.Result{
				"Redis": errors.New("abla habla"),
			},
			args: args{
				"Redis",
			},
			wantKnownCheck: true,
			wantErr:        true,
		},
		{
			name:   "Unknown check",
			result: health.Result{},
			args: args{
				"Redis",
			},
			wantKnownCheck: false,
			wantErr:        false,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotKnownCheck, err := tt.result.CheckResult(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKnownCheck != tt.wantKnownCheck {
				t.Errorf("CheckResult() gotKnownCheck = %v, want %v", gotKnownCheck, tt.wantKnownCheck)
			}
		})
	}
}

func Test_resultWriter_WriteResult(t *testing.T) {
	t.Parallel()
	type args struct {
		checkName string
		result    error
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Successful result",
			args: args{
				checkName: "Sample",
			},
			want: td.Map(health.Result{}, map[interface{}]interface{}{
				"Sample": nil,
			}),
		},
		{
			name: "Error result - simple error",
			args: args{
				checkName: "Sample",
				result:    errors.New("critical error"),
			},
			want: td.Map(health.Result{}, map[interface{}]interface{}{
				"Sample": errors.New("critical error"),
			}),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := health.NewResultWriter()
			r.WriteResult(tt.args.checkName, tt.args.result)
			td.Cmp(t, r.GetResult(), tt.want)
		})
	}
}
