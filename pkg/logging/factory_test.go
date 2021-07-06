package logging

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"go.uber.org/zap"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()
	type args struct {
		levelString string
	}
	type testCase struct {
		name string
		args args
		want zap.AtomicLevel
	}
	tests := []testCase{
		{
			name: "Test parse DEBUG level",
			args: args{
				levelString: "DEBUG",
			},
			want: zap.NewAtomicLevelAt(zap.DebugLevel),
		},
		{
			name: "Test parse DeBuG level",
			args: args{
				levelString: "DeBuG",
			},
			want: zap.NewAtomicLevelAt(zap.DebugLevel),
		},
		{
			name: "Test parse INFO level",
			args: args{
				levelString: "INFO",
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name: "Test parse InFo level",
			args: args{
				levelString: "InFo",
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name: "Test parse WARN level",
			args: args{
				levelString: "WARN",
			},
			want: zap.NewAtomicLevelAt(zap.WarnLevel),
		},
		{
			name: "Test parse WaRn level",
			args: args{
				levelString: "WaRn",
			},
			want: zap.NewAtomicLevelAt(zap.WarnLevel),
		},
		{
			name: "Test parse ERROR level",
			args: args{
				levelString: "ERROR",
			},
			want: zap.NewAtomicLevelAt(zap.ErrorLevel),
		},
		{
			name: "Test parse ErRoR level",
			args: args{
				levelString: "ErRoR",
			},
			want: zap.NewAtomicLevelAt(zap.ErrorLevel),
		},
		{
			name: "Test parse FATAL level",
			args: args{
				levelString: "FATAL",
			},
			want: zap.NewAtomicLevelAt(zap.FatalLevel),
		},
		{
			name: "Test parse FaTaL level",
			args: args{
				levelString: "FaTaL",
			},
			want: zap.NewAtomicLevelAt(zap.FatalLevel),
		},
		{
			name: "Fallback to INFO level if unknown level",
			args: args{
				levelString: "asdf23423",
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name: "Fallback to INFO level if no level",
			args: args{
				levelString: "",
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseLevel(tt.args.levelString)
			td.Cmp(t, got, tt.want)
		})
	}
}

//nolint:paralleltest
func TestConfigureLogging(t *testing.T) {
	type args struct {
		level              zap.AtomicLevel
		developmentLogging bool
		encoding           string
		initialFields      map[string]interface{}
	}
	type testCase struct {
		name string
		args args
	}
	tests := []testCase{
		{
			name: "Test configure defaults",
			args: args{},
		},
		{
			name: "Test configure with initialFields",
			args: args{
				initialFields: map[string]interface{}{
					"asdf": "hello, World",
				},
			},
		},
		{
			name: "Test configure development logging enabled",
			args: args{
				developmentLogging: true,
			},
		},
		{
			name: "Test configure log level",
			args: args{
				level: zap.NewAtomicLevelAt(zap.FatalLevel),
			},
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			ConfigureLogging(tt.args.level, tt.args.developmentLogging, tt.args.encoding, tt.args.initialFields)
			if loggingConfig.Development != tt.args.developmentLogging {
				t.Errorf("loggingConfig.Development = %t, want %t", loggingConfig.Development, tt.args.developmentLogging)
				return
			}

			if loggingConfig.Level != tt.args.level {
				t.Errorf("loggingConfig.Level = %v, want %v", loggingConfig.Level, tt.args.level)
				return
			}

			if tt.args.initialFields != nil {
				td.Cmp(t, loggingConfig.InitialFields, tt.args.initialFields)
			}
		})
	}
}
