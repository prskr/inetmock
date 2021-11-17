package sink_test

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/details"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/wait"
)

var testEvents = []*audit.Event{
	{
		Transport:       auditv1.TransportProtocol_TRANSPORT_PROTOCOL_TCP,
		Application:     auditv1.AppProtocol_APP_PROTOCOL_HTTP,
		SourceIP:        net.ParseIP("127.0.0.1").To4(),
		DestinationIP:   net.ParseIP("127.0.0.1").To4(),
		SourcePort:      32344,
		DestinationPort: 80,
		TLS: &audit.TLSDetails{
			Version:     audit.TLSVersionToEntity(tls.VersionTLS13).String(),
			CipherSuite: tls.CipherSuiteName(tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA),
			ServerName:  "localhost",
		},
		ProtocolDetails: details.HTTP{
			Method: "GET",
			Host:   "localhost",
			URI:    "http://localhost/asdf",
			Proto:  "HTTP 1.1",
			Headers: http.Header{
				"Accept": []string{"application/json"},
			},
		},
	},
	{
		Transport:       auditv1.TransportProtocol_TRANSPORT_PROTOCOL_UDP,
		Application:     auditv1.AppProtocol_APP_PROTOCOL_DNS,
		SourceIP:        net.ParseIP("::1").To16(),
		DestinationIP:   net.ParseIP("::1").To16(),
		SourcePort:      32344,
		DestinationPort: 80,
	},
}

func Test_writerCloserSink_OnSubscribe(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		events []*audit.Event
	}
	tests := []testCase{
		{
			name:   "Get a single event",
			events: testEvents[:1],
		},
		{
			name:   "Get multiple events",
			events: testEvents,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wg := new(sync.WaitGroup)
			wg.Add(len(tt.events))

			ctrl := gomock.NewController(t)

			writerMock := audit_mock.NewMockWriter(ctrl)
			writerMock.
				EXPECT().
				Write(gomock.Any()).
				Do(func(_ *audit.Event) {
					wg.Done()
				}).
				Times(len(tt.events))

			writerCloserSink := sink.NewWriterSink("WriterMock", writerMock)
			var evs audit.EventStream
			var err error

			if evs, err = audit.NewEventStream(logging.CreateTestLogger(t)); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			if err = evs.RegisterSink(ctx, writerCloserSink); err != nil {
				t.Errorf("RegisterSink() error = %v", err)
			}

			for _, ev := range tt.events {
				evs.Emit(*ev)
			}

			select {
			case <-time.After(100 * time.Millisecond):
				t.Errorf("not all events recorded in time")
			case <-wait.ForWaitGroupDone(wg):
			}
		})
	}
}
