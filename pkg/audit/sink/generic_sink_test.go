package sink_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit/sink"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"gitlab.com/inetmock/inetmock/pkg/wait"
)

func Test_genericSink_OnSubscribe(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		events []*audit.Event
	}
	tests := []testCase{
		{
			name:   "Get a single log line",
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

			genericSink := sink.NewGenericSink(t.Name(), func(ev audit.Event) {
				wg.Done()
			})

			var evs audit.EventStream
			var err error
			if evs, err = audit.NewEventStream(logging.CreateTestLogger(t)); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			if err = evs.RegisterSink(ctx, genericSink); err != nil {
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
