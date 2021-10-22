package mock_test

import (
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	mdns "github.com/miekg/dns"

	auditmock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	dnsmock "gitlab.com/inetmock/inetmock/internal/mock/dns"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/protocols/dns"
	"gitlab.com/inetmock/inetmock/protocols/dns/mock"
)

func TestServer_ServeDNS(t *testing.T) {
	t.Parallel()
	type fields struct {
		Handler dns.Handler
	}
	tests := []struct {
		name          string
		fields        fields
		req           *mdns.Msg
		want          interface{}
		wantEmitCalls int
	}{
		{
			name: "Successfully resolve with handler",
			fields: fields{
				Handler: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return &mdns.A{
						A: net.IPv4(10, 10, 0, 1),
					}, nil
				}),
			},
			req: &mdns.Msg{
				Question: []mdns.Question{
					{
						Name:   "www.stackoverflow.com.",
						Qtype:  mdns.TypeA,
						Qclass: mdns.ClassINET,
					},
				},
			},
			want: td.Contains(td.Struct(&mdns.A{
				A: net.IPv4(10, 10, 0, 1),
			}, td.StructFields{})),
			wantEmitCalls: 1,
		},
		{
			name: "Handler does not resolve but returns error",
			fields: fields{
				Handler: dns.HandlerFunc(func(q dns.Question) (dns.ResourceRecord, error) {
					return nil, dns.ErrNoAnswerForQuestion
				}),
			},
			req: &mdns.Msg{
				Question: []mdns.Question{
					{
						Name:   "www.stackoverflow.com.",
						Qtype:  mdns.TypeA,
						Qclass: mdns.ClassINET,
					},
				},
			},
			want:          td.Empty(),
			wantEmitCalls: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			emitter := &auditmock.EmitterMock{
				OnEmit: func(state auditmock.EmitterMockCallsContext, ev audit.Event) {
					if callCount := len(state.Emit); callCount != tt.wantEmitCalls {
						t.Errorf("Want call count: %d, got: %d", tt.wantEmitCalls, callCount)
					}
					t.Logf("Got event: %v", ev)
				},
			}
			s := &mock.Server{
				Name:    t.Name(),
				Handler: tt.fields.Handler,
				Logger:  logging.CreateTestLogger(t),
				Emitter: emitter,
			}

			writerMock := &dnsmock.ResponseWriterMock{
				Local:  new(net.TCPAddr),
				Remote: new(net.TCPAddr),
				OnWriteMsg: func(msg *mdns.Msg) error {
					td.Cmp(t, msg.Answer, tt.want)
					return nil
				},
			}

			s.ServeDNS(writerMock, tt.req)
			td.Cmp(t, len(emitter.Calls.Emit), tt.wantEmitCalls)
		})
	}
}
