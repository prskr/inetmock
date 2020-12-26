package http_mock

import (
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func Test_loadFromConfig(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name        string
		args        args
		wantOptions httpOptions
		wantErr     bool
	}{
		{
			name: "Parse default config",
			args: args{
				config: `
rules:
- pattern: ".*\\.(?i)exe"
  response: ./assets/fakeFiles/sample.exe
`,
			},
			wantOptions: httpOptions{
				Rules: []targetRule{
					{
						pattern: regexp.MustCompile(".*\\.(?i)exe"),
						response: func() string {
							p, _ := filepath.Abs("./assets/fakeFiles/sample.exe")
							return p
						}(),
						requestMatchTarget: RequestMatchTargetPath,
						targetKey:          "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse config with path matcher",
			args: args{
				config: `
rules:
- pattern: ".*\\.(?i)exe"
  matcher: Path 
  response: ./assets/fakeFiles/sample.exe
`,
			},
			wantOptions: httpOptions{
				Rules: []targetRule{
					{
						pattern: regexp.MustCompile(".*\\.(?i)exe"),
						response: func() string {
							p, _ := filepath.Abs("./assets/fakeFiles/sample.exe")
							return p
						}(),
						requestMatchTarget: RequestMatchTargetPath,
						targetKey:          "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse config with header matcher",
			args: args{
				config: `
rules:
- pattern: "^application/octet-stream$"
  target: Content-Type
  matcher: Header
  response: ./assets/fakeFiles/sample.exe
`,
			},
			wantOptions: httpOptions{
				Rules: []targetRule{
					{
						pattern: regexp.MustCompile("^application/octet-stream$"),
						response: func() string {
							p, _ := filepath.Abs("./assets/fakeFiles/sample.exe")
							return p
						}(),
						requestMatchTarget: RequestMatchTargetHeader,
						targetKey:          "Content-Type",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")
			_ = v.ReadConfig(strings.NewReader(tt.args.config))
			gotOptions, err := loadFromConfig(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOptions, tt.wantOptions) {
				t.Errorf("loadFromConfig() gotOptions = %v, want %v", gotOptions, tt.wantOptions)
			}
		})
	}
}
