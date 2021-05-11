package mock_test

import (
	"io/fs"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()
	type fields struct {
		rules        []string
		emitterSetup func(tb testing.TB, ctrl *gomock.Controller) audit.Emitter
		fakeFileFS   fs.FS
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantStatus interface{}
		want       string
	}{
		{
			name: "GET /index.html",
			fields: fields{
				rules: []string{
					`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
				},
				emitterSetup: defaultEmitter,
				fakeFileFS:   defaultFakeFileFS,
			},
			args: args{
				req: tdhttp.NewRequest(http.MethodGet, "https://google.com/index.html", nil),
			},
			want:       defaultHTMLContent,
			wantStatus: td.Between(200, 299),
		},
		{
			name: "GET /profile.htm",
			fields: fields{
				rules: []string{
					`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
				},
				emitterSetup: defaultEmitter,
				fakeFileFS:   defaultFakeFileFS,
			},
			args: args{
				req: tdhttp.NewRequest(http.MethodGet, "https://gitlab.com/profile.htm", nil),
			},
			want:       defaultHTMLContent,
			wantStatus: td.Between(200, 299),
		},
		{
			name: "GET with Accept: text/html",
			fields: fields{
				rules: []string{
					`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
					`Header("Accept", "text/html") => File("default.html")`,
				},
				emitterSetup: defaultEmitter,
				fakeFileFS:   defaultFakeFileFS,
			},
			args: args{
				req: tdhttp.NewRequest(
					http.MethodGet,
					"https://gitlab.com/profile",
					nil,
					http.Header{
						"Accept": []string{"text/html"},
					},
				),
			},
			want:       defaultHTMLContent,
			wantStatus: td.Between(200, 299),
		},
		{
			name: "POST",
			fields: fields{
				rules: []string{
					`METHOD("POST") => Status(204)`,
				},
				emitterSetup: defaultEmitter,
			},
			args: args{
				req: tdhttp.NewRequest(
					http.MethodPost,
					"https://gitlab.com/profile",
					nil,
					http.Header{
						"Accept": []string{"text/html"},
					},
				),
			},
			want:       "",
			wantStatus: 204,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			ctrl := gomock.NewController(t)

			router := &mock.Router{
				HandlerName: t.Name(),
				Logger:      logger,
				Emitter:     tt.fields.emitterSetup(t, ctrl),
				FakeFileFS:  tt.fields.fakeFileFS,
			}

			for _, rule := range tt.fields.rules {
				if err := router.RegisterRule(rule); !td.CmpNoError(t, err) {
					return
				}
			}

			tdhttp.NewTestAPI(t, router).
				Request(tt.args.req).
				CmpStatus(tt.wantStatus).
				CmpBody(tt.want)
		})
	}
}

func defaultEmitter(tb testing.TB, ctrl *gomock.Controller) audit.Emitter {
	tb.Helper()
	emitter := audit_mock.NewMockEmitter(ctrl)
	emitter.
		EXPECT().
		Emit(gomock.Any()).
		Times(1)
	return emitter
}
