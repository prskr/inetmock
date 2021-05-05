package mock_test

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
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
				req: &http.Request{
					URL:    mustParseURL("https://google.com/index.html"),
					Method: http.MethodGet,
				},
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
				req: &http.Request{
					URL:    mustParseURL("https://gitlab.com/profile.htm"),
					Method: http.MethodGet,
				},
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
				req: &http.Request{
					URL: mustParseURL("https://gitlab.com/profile"),
					Header: http.Header{
						"Accept": []string{"text/html"},
					},
					Method: http.MethodGet,
				},
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
				req: &http.Request{
					URL:    mustParseURL("https://gitlab.com/profile"),
					Header: http.Header{},
					Method: http.MethodPost,
				},
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

			client := setupHTTPServer(t, router)

			resp, err := client.Do(tt.args.req)
			if tt.wantErr == td.CmpNoError(t, err) {
				return
			}

			t.Cleanup(func() {
				td.CmpNoError(t, resp.Body.Close())
			})

			td.Cmp(t, resp.StatusCode, tt.wantStatus)

			builder := new(strings.Builder)
			_, err = io.Copy(builder, resp.Body)
			td.CmpNoError(t, err)
			td.Cmp(t, builder.String(), tt.want)
		})
	}
}

func mustParseURL(urlToParse string) *url.URL {
	parsed, err := url.Parse(urlToParse)
	if err != nil {
		panic(err)
	}
	return parsed
}

func setupHTTPServer(tb testing.TB, handler http.Handler) *http.Client {
	tb.Helper()
	listener := eptest.NewInMemoryListener(tb)

	go func() {
		switch err := http.Serve(listener, handler); {
		case errors.Is(err, nil), errors.Is(err, http.ErrServerClosed):
		default:
			tb.Errorf("http.Serve() error = %v", err)
		}
	}()

	return eptest.HTTPClientForInMemListener(listener)
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
