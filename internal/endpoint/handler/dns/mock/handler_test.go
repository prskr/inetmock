package mock_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.com/inetmock/inetmock/internal/endpoint"
	"gitlab.com/inetmock/inetmock/internal/endpoint/eptest"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/dns/mock"
	auditmock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/internal/test"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func Test_dnsHandler_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		opts map[string]interface{}
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    []net.IP
		wantErr bool
	}{
		{
			name: "Resolve all to 1.1.1.1",
			args: args{
				opts: map[string]interface{}{
					"ttl": "30s",
					"cache": map[string]interface{}{
						"type": "none",
					},
					"rules": []string{
						`=> IP(1.1.1.1)`,
					},
					"default": map[string]interface{}{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
				},
				host: "google.com",
			},
			want: []net.IP{
				net.IPv4(1, 1, 1, 1),
			},
			wantErr: false,
		},
		{
			name: "Resolve with fallback",
			args: args{
				opts: map[string]interface{}{
					"ttl": "30s",
					"cache": map[string]interface{}{
						"type": "none",
					},
					"rules": make([]string, 0),
					"default": map[string]interface{}{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
				},
				host: "google.com",
			},
			want: []net.IP{
				net.IPv4(10, 10, 0, 1),
			},
			wantErr: false,
		},
		{
			name: "Resolve google.com domain",
			args: args{
				opts: map[string]interface{}{
					"ttl": "30s",
					"cache": map[string]interface{}{
						"type": "none",
					},
					"rules": []string{
						`A(".*\\.google\\.com\\.$") => IP(1.1.1.1)`,
					},
					"default": map[string]interface{}{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
				},
				host: "mail.google.com",
			},
			want: []net.IP{
				net.IPv4(1, 1, 1, 1),
			},
			wantErr: false,
		},
		{
			name: "Resolve stackoverflow.com domain",
			args: args{
				opts: map[string]interface{}{
					"ttl": "30s",
					"cache": map[string]interface{}{
						"type": "none",
					},
					"rules": []string{
						`A(".*\\.google\\.com\\.$") => IP(1.1.1.1)`,
						`A(".*\\.stackoverflow\\.com\\.$") => IP(1.2.3.4)`,
					},
					"default": map[string]interface{}{
						"type": "incremental",
						"cidr": "10.10.0.0/16",
					},
				},
				host: "www.stackoverflow.com",
			},
			want: []net.IP{
				net.IPv4(1, 2, 3, 4),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			listener := eptest.NewInMemoryListener(t)
			ctx, cancel := context.WithCancel(test.Context(t))
			t.Cleanup(cancel)
			emitter := auditmock.NewMockEmitter(ctrl)
			lifecycle := endpoint.NewEndpointLifecycle(t.Name(), endpoint.Uplink{Listener: listener}, tt.args.opts)
			if !tt.wantErr {
				emitter.EXPECT().
					Emit(gomock.Any()).
					MinTimes(1)
			}
			handler := mock.New(logging.CreateTestLogger(t), emitter)
			if err := handler.Start(ctx, lifecycle); err != nil {
				if !tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			resolver := eptest.DNSResolverForInMemListener(listener)
			requestCtx, requestCancel := context.WithTimeout(ctx, 50*time.Millisecond)
			t.Cleanup(requestCancel)
			if ips, err := resolver.LookupIP(requestCtx, "ip", tt.args.host); err != nil {
				if !tt.wantErr {
					t.Errorf("LookupIP() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else {
				for _, gotIP := range ips {
					var matched bool
					for _, wantIP := range tt.want {
						matched = matched || wantIP.Equal(gotIP)
					}
					if !matched {
						t.Errorf("Got %v but didn't expect it", gotIP)
					}
				}
			}
		})
	}
}
