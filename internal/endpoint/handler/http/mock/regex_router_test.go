package mock_test

import (
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	defaultHTMLContent = `<html>
<head>
    <title>INetSim default HTML page</title>
</head>
<body>
<p></p>
<p align="center">This is the default HTML page for INetMock HTTP mock handler.</p>
<p align="center">This file is an HTML document.</p>
</body>
</html>`
)

//nolint:funlen
func TestRegexpHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	if err := mock.InitMetrics(); !td.CmpNoError(t, err) {
		return
	}
	type fields struct {
		rules        []mock.TargetRule
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
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, "default.html"),
				},
				emitterSetup: defaultEmitter,
				fakeFileFS: fstest.MapFS{
					"default.html": &fstest.MapFile{
						Data:    []byte(defaultHTMLContent),
						ModTime: time.Now().Add(-1337 * time.Second),
					},
				},
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
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, "default.html"),
				},
				emitterSetup: defaultEmitter,
				fakeFileFS: fstest.MapFS{
					"default.html": &fstest.MapFile{
						Data:    []byte(defaultHTMLContent),
						ModTime: time.Now().Add(-1337 * time.Second),
					},
				},
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
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, "default.html"),
					mock.MustHeaderTargetRule("Accept", "(?i)text/html", "default.html"),
				},
				emitterSetup: defaultEmitter,
				fakeFileFS: fstest.MapFS{
					"default.html": &fstest.MapFile{
						Data:    []byte(defaultHTMLContent),
						ModTime: time.Now().Add(-1337 * time.Second),
					},
				},
			},
			args: args{
				req: &http.Request{
					URL: mustParseURL("https://gitlab.com/profile"),
					Header: map[string][]string{
						"Accept": {"text/html"},
					},
					Method: http.MethodGet,
				},
			},
			want:       defaultHTMLContent,
			wantStatus: td.Between(200, 299),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			ctrl := gomock.NewController(t)
			h := mock.NewRegexHandler(t.Name(), logger, tt.fields.emitterSetup(t, ctrl), tt.fields.fakeFileFS)

			for _, rule := range tt.fields.rules {
				h.AddRouteRule(rule)
			}

			client := setupHTTPServer(t, h)

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
	srv := &http.Server{
		Handler: handler,
	}
	listener := eptest.NewInMemoryListener(tb)

	go func() {
		td.Cmp(tb, srv.Serve(listener), http.ErrServerClosed)
	}()

	tb.Cleanup(func() {
		td.CmpNoError(tb, srv.Close())
	})

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
