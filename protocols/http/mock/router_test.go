package mock_test

import (
	"io/fs"
	"net/http"
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/http/mock"
)

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()
	type fields struct {
		rules      []string
		fakeFileFS fs.FS
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantStatus any
		want       string
	}{
		{
			name: "GET /index.html",
			fields: fields{
				rules: []string{
					`PathPattern("\\.(?i)(htm|html)$") => File("default.html")`,
				},
				fakeFileFS: defaultFakeFileFS,
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
				fakeFileFS: defaultFakeFileFS,
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
				fakeFileFS: defaultFakeFileFS,
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

			router := &mock.Router{
				HandlerName: t.Name(),
				Logger:      logger,
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
