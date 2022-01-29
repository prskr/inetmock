package mock_test

import (
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	httpmock "gitlab.com/inetmock/inetmock/internal/mock/http"
	"gitlab.com/inetmock/inetmock/internal/rules"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestStatusHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		args                []rules.Param
		responseWriterSetup func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter
		request             *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get status 204",
			args: args{
				args: []rules.Param{
					{
						Int: rules.IntP(204),
					},
				},
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().WriteHeader(204)
					return rwMock
				},
				request: new(http.Request),
			},
			wantErr: false,
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
			ctrl := gomock.NewController(t)
			got, err := mock.StatusHandler(logger, nil, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StatusHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got.ServeHTTP(tt.args.responseWriterSetup(t, ctrl), tt.args.request)
		})
	}
}

func TestFileHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		fakeFileFS          fs.FS
		args                []rules.Param
		responseWriterSetup func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter
		request             *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Header().Return(http.Header{}).MinTimes(1)
					rwMock.EXPECT().WriteHeader(200)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.Contains("<title>INetSim default HTML page</title>")))
					return rwMock
				},
				request: new(http.Request),
			},
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Header().Return(http.Header{}).MinTimes(1)
					rwMock.EXPECT().WriteHeader(500)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.Contains("open default.html: file does not exist")))
					return rwMock
				},
				request: new(http.Request),
			},
			wantErr: false,
		},
		{
			name: "Expect error due to missing argument",
			args: args{
				fakeFileFS: nil,
				args:       []rules.Param{},
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					return rwMock
				},
				request: new(http.Request),
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					return rwMock
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
			ctrl := gomock.NewController(t)
			got, err := mock.FileHandler(logger, tt.args.fakeFileFS, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("StatusHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got.ServeHTTP(tt.args.responseWriterSetup(t, ctrl), tt.args.request)
		})
	}
}

func TestJSONHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		args                []rules.Param
		responseWriterSetup func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter
		request             *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    http.Handler
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.String(`{}`)))
					return rwMock
				},
				request: new(http.Request),
			},
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.String(`{"Name": "Ted Tester"}`)))
					return rwMock
				},
				request: new(http.Request),
			},
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.String(`{"Name": "Ted Tester", "Address": {"Street": "Some street 1"}}`)))
					return rwMock
				},
				request: new(http.Request),
			},
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.String(`{"Name": "Ted Tester", "Colleagues": []}`)))
					return rwMock
				},
				request: new(http.Request),
			},
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
				responseWriterSetup: func(tb testing.TB, ctrl *gomock.Controller) http.ResponseWriter {
					tb.Helper()
					rwMock := httpmock.NewMockResponseWriter(ctrl)
					rwMock.EXPECT().Write(test.GenericMatcher(tb, td.String(`{"Name": "Ted Tester", "Colleagues": [{"Name": "Carl"}]}`)))
					return rwMock
				},
				request: new(http.Request),
			},
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
			ctrl := gomock.NewController(t)
			got, err := mock.JSONHandler(logger, nil, tt.args.args...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("JSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got.ServeHTTP(tt.args.responseWriterSetup(t, ctrl), tt.args.request)
		})
	}
}
