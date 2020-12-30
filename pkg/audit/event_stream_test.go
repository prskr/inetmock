package audit_test

import (
	"net"
	"sync"
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

var (
	defaultSink = &testSink{
		name: "test defaultSink",
	}
)

type testSink struct {
	name     string
	consumer func(event audit.Event)
}

func (t *testSink) Name() string {
	return t.name
}

func (t *testSink) OnSubscribe(evs <-chan audit.Event) {
	go func() {
		for ev := range evs {
			if t.consumer != nil {
				t.consumer(ev)
			}
		}
	}()
}

func wgMockSink(t testing.TB, wg *sync.WaitGroup) audit.Sink {
	return &testSink{
		name: "WG mock sink",
		consumer: func(event audit.Event) {
			t.Logf("Got event = %v", event)
			wg.Done()
		},
	}
}

func Test_eventStream_RegisterSink(t *testing.T) {
	type args struct {
		s audit.Sink
	}

	tests := []struct {
		name    string
		args    args
		setup   func(e audit.EventStream)
		wantErr bool
	}{
		{
			name: "Register test defaultSink",
			args: args{
				s: defaultSink,
			},
			wantErr: false,
		},
		{
			name: "Fail due to already registered defaultSink",
			args: args{
				s: defaultSink,
			},
			setup: func(e audit.EventStream) {
				_ = e.RegisterSink(defaultSink)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var e audit.EventStream
			if e, err = audit.NewEventStream(logging.CreateTestLogger(t)); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			t.Cleanup(func() {
				_ = e.Close()
			})

			if tt.setup != nil {
				tt.setup(e)
			}

			if err := e.RegisterSink(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("RegisterSink() error = %v, wantErr %v", err, tt.wantErr)
			}

			found := false
			for _, s := range e.Sinks() {
				if found = s == tt.args.s.Name(); found {
					break
				}
			}
			if !found {
				t.Errorf("expected defaultSink name %s not found in registered sinks %v", tt.args.s.Name(), e.Sinks())
			}
		})
	}
}

func Test_eventStream_Emit(t *testing.T) {
	type args struct {
		evs  []audit.Event
		opts []audit.EventStreamOption
	}
	tests := []struct {
		name      string
		args      args
		subscribe bool
	}{
		{
			name:      "Expect to get a single event",
			subscribe: true,
			args: args{
				opts: []audit.EventStreamOption{audit.WithBufferSize(10)},
				evs: []audit.Event{
					{
						Transport:       audit.TransportProtocol_TCP,
						Application:     audit.AppProtocol_HTTP,
						SourceIP:        net.ParseIP("127.0.0.1"),
						DestinationIP:   net.ParseIP("127.0.0.1"),
						SourcePort:      32344,
						DestinationPort: 80,
					},
				},
			},
		},
		{
			name:      "Expect to get multiple events",
			subscribe: true,
			args: args{
				opts: []audit.EventStreamOption{audit.WithBufferSize(10)},
				evs: []audit.Event{
					{
						Transport:       audit.TransportProtocol_TCP,
						Application:     audit.AppProtocol_HTTP,
						SourceIP:        net.ParseIP("127.0.0.1"),
						DestinationIP:   net.ParseIP("127.0.0.1"),
						SourcePort:      32344,
						DestinationPort: 80,
					},
					{
						Transport:       audit.TransportProtocol_TCP,
						Application:     audit.AppProtocol_DNS,
						SourceIP:        net.ParseIP("::1"),
						DestinationIP:   net.ParseIP("::1"),
						SourcePort:      32344,
						DestinationPort: 80,
					},
				},
			},
		},
		{
			name: "Emit without subscribe sink",
			args: args{
				opts: []audit.EventStreamOption{audit.WithBufferSize(0)},
				evs: []audit.Event{
					{
						Transport:       audit.TransportProtocol_TCP,
						Application:     audit.AppProtocol_HTTP,
						SourceIP:        net.ParseIP("127.0.0.1"),
						DestinationIP:   net.ParseIP("127.0.0.1"),
						SourcePort:      32344,
						DestinationPort: 80,
					},
				},
			},
			subscribe: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var e audit.EventStream
			if e, err = audit.NewEventStream(logging.CreateTestLogger(t), tt.args.opts...); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			t.Cleanup(func() {
				_ = e.Close()
			})

			emittedWaitGroup := &sync.WaitGroup{}
			receivedWaitGroup := &sync.WaitGroup{}

			emittedWaitGroup.Add(len(tt.args.evs))

			if tt.subscribe {
				receivedWaitGroup.Add(len(tt.args.evs))
				if err := e.RegisterSink(wgMockSink(t, receivedWaitGroup)); err != nil {
					t.Errorf("RegisterSink() error = %v", err)
				}
			}

			go func(evs []audit.Event, wg *sync.WaitGroup) {
				for _, ev := range evs {
					e.Emit(ev)
					wg.Done()
				}
			}(tt.args.evs, emittedWaitGroup)

			select {
			case <-waitGroupDone(emittedWaitGroup):
			case <-time.After(100 * time.Millisecond):
				t.Errorf("not all events emitted in time")
			}

			select {
			case <-waitGroupDone(receivedWaitGroup):
			case <-time.After(5 * time.Second):
				t.Errorf("did not get all expected events in time")
			}
		})
	}
}

func waitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(done)
	}(wg)

	return done
}
