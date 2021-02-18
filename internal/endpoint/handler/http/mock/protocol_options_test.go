package mock

import (
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mitchellh/mapstructure"

	endpoint_mock "gitlab.com/inetmock/inetmock/internal/mock/endpoint"
)

//nolint:funlen
func Test_loadFromConfig(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	type testCase struct {
		name string
		args struct {
			config map[string]interface{}
		}
		wantOptions httpOptions
		wantErr     bool
	}
	tests := []testCase{
		{
			name: "Parse default config",
			args: args{
				config: map[string]interface{}{
					"rules": []struct {
						Pattern  string
						Matcher  string
						Response string
					}{
						{
							Pattern:  ".*\\.(?i)exe",
							Response: "./assets/fakeFiles/sample.exe",
						},
					},
				},
			},
			wantOptions: httpOptions{
				Rules: []targetRule{
					{
						pattern: regexp.MustCompile(`.*\.(?i)exe`),
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
				config: map[string]interface{}{
					"rules": []struct {
						Pattern  string
						Matcher  string
						Response string
					}{
						{
							Pattern:  ".*\\.(?i)exe",
							Response: "./assets/fakeFiles/sample.exe",
							Matcher:  "Path",
						},
					},
				},
			},
			wantOptions: httpOptions{
				Rules: []targetRule{
					{
						pattern: regexp.MustCompile(`.*\.(?i)exe`),
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
				config: map[string]interface{}{
					"rules": []struct {
						Pattern  string
						Matcher  string
						Target   string
						Response string
					}{
						{
							Pattern:  "^application/octet-stream$",
							Response: "./assets/fakeFiles/sample.exe",
							Target:   "Content-Type",
							Matcher:  "Header",
						},
					},
				},
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
		{
			name: "Parse config with header matcher and TLS true",
			args: args{
				config: map[string]interface{}{
					"tls": true,
					"rules": []struct {
						Pattern  string
						Matcher  string
						Target   string
						Response string
					}{
						{
							Pattern:  "^application/octet-stream$",
							Response: "./assets/fakeFiles/sample.exe",
							Target:   "Content-Type",
							Matcher:  "Header",
						},
					},
				},
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			lcMock := endpoint_mock.NewMockLifecycle(ctrl)

			lcMock.EXPECT().UnmarshalOptions(gomock.Any()).Do(func(cfg interface{}) {
				_ = mapstructure.Decode(tt.args.config, cfg)
			})

			gotOptions, err := loadFromConfig(lcMock)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOptions, tt.wantOptions) {
				t.Errorf("loadFromConfig() gotOptions = %v, want %v", gotOptions, tt.wantOptions)
			}
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
