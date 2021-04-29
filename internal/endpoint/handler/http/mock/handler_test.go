//nolint:funlen
package mock_test

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	v1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func Test_httpHandler_Start(t *testing.T) {
	t.Parallel()
	if !td.CmpNoError(t, mock.InitMetrics()) {
		return
	}
	type fields struct {
		fakeFileFS fs.FS
	}
	type args struct {
		opts map[string]interface{}
		req  *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus interface{}
		wantBody   interface{}
		wantEvent  interface{}
		wantErr    bool
	}{
		{
			name: "Get /index.html",
			fields: fields{
				fakeFileFS: defaultFakeFileFS,
			},
			args: args{
				opts: map[string]interface{}{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: mustParseURL("https://www.google.de/index.html"),
				},
			},
			wantEvent: td.Struct(audit.Event{}, td.StructFields{
				"Application": v1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(details.HTTP{}, td.StructFields{
					"Host":   "www.google.de",
					"URI":    "/index.html",
					"Method": http.MethodGet,
				}),
			}),
			wantStatus: 200,
			wantBody:   defaultHTMLContent,
			wantErr:    false,
		},
		{
			name: "Get /asdf.html",
			fields: fields{
				fakeFileFS: defaultFakeFileFS,
			},
			args: args{
				opts: map[string]interface{}{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: mustParseURL("https://www.google.de/asdf.html"),
				},
			},
			wantEvent: td.Struct(audit.Event{}, td.StructFields{
				"Application": v1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(details.HTTP{}, td.StructFields{
					"Host":   "www.google.de",
					"URI":    "/asdf.html",
					"Method": http.MethodGet,
				}),
			}),
			wantStatus: 200,
			wantBody:   defaultHTMLContent,
			wantErr:    false,
		},
		{
			name: "Get /asdf with HTML accept header",
			fields: fields{
				fakeFileFS: defaultFakeFileFS,
			},
			args: args{
				opts: map[string]interface{}{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
						`Header("Accept", "text/html") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: mustParseURL("https://www.google.de/asdf"),
					Header: http.Header{
						"Accept": []string{"text/html"},
					},
				},
			},
			wantEvent: td.Struct(audit.Event{}, td.StructFields{
				"Application": v1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(details.HTTP{}, td.StructFields{
					"Host":   "www.google.de",
					"URI":    "/asdf",
					"Method": http.MethodGet,
				}),
			}),
			wantStatus: 200,
			wantBody:   defaultHTMLContent,
			wantErr:    false,
		},
		{
			name: "Error because of syntax error in rule",
			args: args{
				opts: map[string]interface{}{
					"rules": []string{
						`> File("default.html")`,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			logger := logging.CreateTestLogger(t)
			listener := eptest.NewInMemoryListener(t)
			lifecycle := endpoint.NewEndpointLifecycle(t.Name(), endpoint.Uplink{Listener: listener}, tt.args.opts)
			emitterMock := audit_mock.NewMockEmitter(ctrl)
			if !tt.wantErr {
				emitterMock.EXPECT().Emit(test.GenericMatcher(t, tt.wantEvent))
			}
			handler := mock.New(logger, emitterMock, tt.fields.fakeFileFS)
			if err := handler.Start(ctx, lifecycle); err != nil {
				if !tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			client := eptest.HTTPClientForInMemListener(listener)
			if resp, err := client.Do(tt.args.req); err != nil {
				if !tt.wantErr {
					t.Errorf("client.Do() error = %v", err)
				}
				return
			} else {
				if !td.Cmp(t, resp.StatusCode, tt.wantStatus) {
					return
				}
				t.Cleanup(func() {
					_ = resp.Body.Close()
				})
				var bodyBuilder = new(strings.Builder)
				if _, err := io.Copy(bodyBuilder, resp.Body); !td.CmpNoError(t, err) {
					return
				}
				td.Cmp(t, bodyBuilder.String(), tt.wantBody)
			}
		})
	}
}
