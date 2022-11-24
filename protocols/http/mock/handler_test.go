package mock_test

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/internal/endpoint"
	audit_mock "inetmock.icb4dc0.de/inetmock/internal/mock/audit"
	"inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/protocols/http/mock"
)

func Test_httpHandler_Start(t *testing.T) {
	t.Parallel()
	type fields struct {
		fakeFileFS fs.FS
	}
	type args struct {
		opts map[string]any
		req  *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus any
		wantBody   any
		wantEvent  any
		wantErr    bool
	}{
		{
			name: "Get /index.html",
			fields: fields{
				fakeFileFS: defaultFakeFileFS,
			},
			args: args{
				opts: map[string]any{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: test.MustParseURL("https://www.google.de/index.html"),
				},
			},
			wantEvent: td.Struct(new(audit.Event), td.StructFields{
				"Application": auditv1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(new(audit.HTTP), td.StructFields{
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
				opts: map[string]any{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: test.MustParseURL("https://www.google.de/asdf.html"),
				},
			},
			wantEvent: td.Struct(new(audit.Event), td.StructFields{
				"Application": auditv1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(new(audit.HTTP), td.StructFields{
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
				opts: map[string]any{
					"rules": []string{
						`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
						`Header("Accept", "text/html") => File("default.html")`,
					},
				},
				req: &http.Request{
					URL: test.MustParseURL("https://www.google.de/asdf"),
					Header: http.Header{
						"Accept": []string{"text/html"},
					},
				},
			},
			wantEvent: td.Struct(new(audit.Event), td.StructFields{
				"Application": auditv1.AppProtocol_APP_PROTOCOL_HTTP,
				"ProtocolDetails": td.Struct(new(audit.HTTP), td.StructFields{
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
				opts: map[string]any{
					"rules": []string{
						`= > File("default.html")`,
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
			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			logger := logging.CreateTestLogger(t)
			listener := test.NewInMemoryListener(t)
			lifecycle := endpoint.NewStartupSpec(t.Name(), endpoint.NewUplink(listener), tt.args.opts)
			emitterMock := new(audit_mock.EmitterMock)

			if !tt.wantErr {
				t.Cleanup(func() {
					emitterMock.WithCalls(func(calls *audit_mock.EmitterMockCalls) {
						for _, call := range calls.Emit() {
							td.Cmp(t, call.Params.Ev, tt.wantEvent)
						}
					})
				})
			}
			handler := mock.New(logger, emitterMock, tt.fields.fakeFileFS)
			if err := handler.Start(ctx, lifecycle); err != nil {
				if !tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			client := test.HTTPClientForInMemListener(listener)
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
				bodyBuilder := new(strings.Builder)
				if _, err := io.Copy(bodyBuilder, resp.Body); !td.CmpNoError(t, err) {
					return
				}
				td.Cmp(t, bodyBuilder.String(), tt.wantBody)
			}
		})
	}
}
