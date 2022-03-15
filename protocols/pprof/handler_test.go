package pprof_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/pprof"
)

func Test_pprofHandler_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus any
		wantEvent  any
	}{
		{
			name: "Expect /debug/pprof/ index to succeed",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/debug/pprof/"),
				},
			},
			wantErr:    false,
			wantStatus: 200,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
		{
			name: "Expect /debug/pprof/cmdline call to succeed",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/debug/pprof/cmdline?seconds=1"),
				},
			},
			wantErr:    false,
			wantStatus: 200,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
		{
			name: "Expect /debug/pprof/profile call to succeed",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/debug/pprof/profile?seconds=1"),
				},
			},
			wantErr:    false,
			wantStatus: 200,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
		{
			name: "Expect /debug/pprof/symbol call to succeed",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/debug/pprof/symbol?seconds=1"),
				},
			},
			wantErr:    false,
			wantStatus: 200,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
		{
			name: "Expect /debug/pprof/trace call to succeed",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/debug/pprof/trace?seconds=1"),
				},
			},
			wantErr:    false,
			wantStatus: 200,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
		{
			name: "Expect / to return 404",
			args: args{
				&http.Request{
					URL: mustParseURL("http://localhost/"),
				},
			},
			wantErr:    false,
			wantStatus: 404,
			wantEvent:  td.Struct(audit.Event{}, td.StructFields{}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			emitter := audit_mock.NewMockEmitter(ctrl)
			if !tt.wantErr {
				emitter.EXPECT().Emit(test.GenericMatcher(t, tt.wantEvent)).MinTimes(1)
			}
			p := pprof.New(logging.CreateTestLogger(t), emitter)

			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			listener := test.NewInMemoryListener(t)
			lifecycle := endpoint.NewStartupSpec(t.Name(), endpoint.NewUplink(listener), nil)

			if err := p.Start(ctx, lifecycle); err != nil {
				if !tt.wantErr {
					t.Errorf("pprofHandler.Start() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			client := test.HTTPClientForInMemListener(listener)

			if resp, err := client.Do(tt.args.req); err != nil {
				if !tt.wantErr {
					t.Errorf("client.Do() error = %v", err)
				}
			} else {
				if !td.Cmp(t, resp.StatusCode, tt.wantStatus) {
					return
				}
				t.Cleanup(func() {
					_ = resp.Body.Close()
				})
			}
		})
	}
}

func mustParseURL(rawURL string) *url.URL {
	if u, err := url.Parse(rawURL); err != nil {
		panic(err)
	} else {
		return u
	}
}
