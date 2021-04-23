package mock_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
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
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantStatus   interface{}
		wantRespHash string
	}{
		{
			name: "GET /index.html",
			fields: fields{
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, filepath.Join("testdata", "default.html")),
				},
				emitterSetup: defaultEmitter,
			},
			args: args{
				req: &http.Request{
					URL:    mustParseURL("https://google.com/index.html"),
					Method: http.MethodGet,
				},
			},
			wantRespHash: "c2a3f8995831dd1e79cb753619a55752692168f6cf846b07405f2070492f481c",
			wantStatus:   td.Between(200, 299),
		},
		{
			name: "GET /profile.htm",
			fields: fields{
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, filepath.Join("testdata", "default.html")),
				},
				emitterSetup: defaultEmitter,
			},
			args: args{
				req: &http.Request{
					URL:    mustParseURL("https://gitlab.com/profile.htm"),
					Method: http.MethodGet,
				},
			},
			wantRespHash: "c2a3f8995831dd1e79cb753619a55752692168f6cf846b07405f2070492f481c",
			wantStatus:   td.Between(200, 299),
		},
		{
			name: "GET with Accept: text/html",
			fields: fields{
				rules: []mock.TargetRule{
					mock.MustPathTargetRule(`\.(?i)(htm|html)$`, filepath.Join("testdata", "default.html")),
					mock.MustHeaderTargetRule("Accept", "(?i)text/html", filepath.Join("testdata", "default.html")),
				},
				emitterSetup: defaultEmitter,
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
			wantRespHash: "c2a3f8995831dd1e79cb753619a55752692168f6cf846b07405f2070492f481c",
			wantStatus:   td.Between(200, 299),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger := logging.CreateTestLogger(t)
			ctrl := gomock.NewController(t)
			h := mock.NewRegexHandler(t.Name(), logger, tt.fields.emitterSetup(t, ctrl))

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

			sha256Hash := sha256.New()
			_, err = io.Copy(sha256Hash, resp.Body)
			td.CmpNoError(t, err)
			computedHash := hex.EncodeToString(sha256Hash.Sum(nil))
			td.Cmp(t, computedHash, tt.wantRespHash)
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
