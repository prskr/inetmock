package mock_test

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestStatusHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		args    []rules.Param
		request *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "Get status 204",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(http.StatusNoContent),
					},
				},
				request: new(http.Request),
			},
			wantErr: false,
			want: td.Struct(&http.Response{
				StatusCode: http.StatusNoContent,
			}, td.StructFields{}),
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				args:    []rules.Param{},
				request: new(http.Request),
			},
			wantErr: true,
		},
		{
			name: "Expect error due to argument type mismatch",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP("Hello, World"),
					},
				},
				request: new(http.Request),
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			handler, err := mock.StatusHandler(logger, nil, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StatusHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, tt.args.request)
			got := recorder.Result()
			td.Cmp(t, got, tt.want)
		})
	}
}

func TestFileHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		fakeFileFS fs.FS
		args       []rules.Param
		request    *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "Get default HTML content",
			args: args{
				fakeFileFS: defaultFakeFileFS,
				args: []rules.Param{
					{
						String: rules.StringP("default.html"),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.Contains("<title>INetSim default HTML page</title>")),
			}),
			wantErr: false,
		},
		{
			name: "Expect 500 due to error in FS",
			args: args{
				fakeFileFS: fstest.MapFS{},
				args: []rules.Param{
					{
						String: rules.StringP("default.html"),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusInternalServerError,
			}, td.StructFields{
				"Body": readerSmuggle(td.Contains("open default.html: file does not exist")),
			}),
			wantErr: false,
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				fakeFileFS: nil,
				args:       []rules.Param{},
				request:    new(http.Request),
			},
			wantErr: true,
		},
		{
			name: "Expect error due to argument type mismatch",
			args: args{
				fakeFileFS: nil,
				args: []rules.Param{
					{
						Int: rules.IntP(42),
					},
				},
				request: new(http.Request),
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			fileHandler, err := mock.FileHandler(logger, tt.args.fakeFileFS, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StatusHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			recorder := httptest.NewRecorder()
			fileHandler.ServeHTTP(recorder, tt.args.request)
			td.Cmp(t, recorder.Result(), tt.want)
		})
	}
}

func TestJSONHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		args    []rules.Param
		request *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Get status empty JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{}`),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.String(`{}`)),
			}),
			wantErr: false,
		},
		{
			name: "Get non-empty JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{"Name": "Ted Tester"}`),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.String(`{"Name": "Ted Tester"}`)),
			}),
			wantErr: false,
		},
		{
			name: "Get nested JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{"Name": "Ted Tester", "Address": {"Street": "Some street 1"}}`),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.String(`{"Name": "Ted Tester", "Address": {"Street": "Some street 1"}}`)),
			}),
			wantErr: false,
		},
		{
			name: "Get nested empty array JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{"Name": "Ted Tester", "Colleagues": []}`),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.String(`{"Name": "Ted Tester", "Colleagues": []}`)),
			}),
			wantErr: false,
		},
		{
			name: "Get nested array JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{"Name": "Ted Tester", "Colleagues": [{"Name": "Carl"}]}`),
					},
				},
				request: new(http.Request),
			},
			want: td.Struct(&http.Response{
				StatusCode: http.StatusOK,
			}, td.StructFields{
				"Body": readerSmuggle(td.String(`{"Name": "Ted Tester", "Colleagues": [{"Name": "Carl"}]}`)),
			}),
			wantErr: false,
		},
		{
			name: "Invalid JSON",
			args: args{
				args: []rules.Param{
					{
						String: rules.StringP(`{`),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			got, err := mock.JSONHandler(logger, nil, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("JSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			recorder := httptest.NewRecorder()
			got.ServeHTTP(recorder, tt.args.request)
			td.Cmp(t, recorder.Result(), tt.want)
		})
	}
}

func readerSmuggle(expected interface{}) interface{} {
	return td.Smuggle(func(reader io.Reader) (string, error) {
		data, err := io.ReadAll(reader)
		return string(data), err
	}, expected)
}
