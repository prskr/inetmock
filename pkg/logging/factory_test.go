package logging

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

//nolint:funlen
func TestParseLevel(t *testing.T) {
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			if got := ParseLevel(tt.args.levelString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.want)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}

func TestConfigureLogging(t *testing.T) {
	type args struct {
		level              zap.AtomicLevel
		developmentLogging bool
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			ConfigureLogging(tt.args.level, tt.args.developmentLogging, tt.args.initialFields)
			if loggingConfig.Development != tt.args.developmentLogging {
				t.Errorf("loggingConfig.Development = %t, want %t", loggingConfig.Development, tt.args.developmentLogging)
				return
			}

			if loggingConfig.Level != tt.args.level {
				t.Errorf("loggingConfig.Level = %v, want %v", loggingConfig.Level, tt.args.level)
				return
			}

			if tt.args.initialFields != nil && !reflect.DeepEqual(loggingConfig.InitialFields, tt.args.initialFields) {
				t.Errorf("loggingConfig.InitialFields = %v, want %v", loggingConfig.InitialFields, tt.args.initialFields)
				return
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
