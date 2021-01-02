package audit_test

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/audit"
	"gitlab.com/inetmock/inetmock/pkg/logging"
)

func Test_writerCloserSink_OnSubscribe(t *testing.T) {
	type testCase struct {
		name   string
		events []audit.Event
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
	scenario := func(tt testCase) func(t *testing.T) {
		return func(t *testing.T) {
			wg := new(sync.WaitGroup)
			wg.Add(len(tt.events))

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			writerMock := audit_mock.NewMockWriter(ctrl)
			writerMock.
				EXPECT().
				Write(gomock.Any()).
				Do(func(_ *audit.Event) {
					wg.Done()
				}).
				Times(len(tt.events))

			writerCloserSink := audit.NewWriterSink("WriterMock", writerMock, audit.WithCloseOnExit)
			var evs audit.EventStream
			var err error

			if evs, err = audit.NewEventStream(logging.CreateTestLogger(t)); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			if err = evs.RegisterSink(writerCloserSink); err != nil {
				t.Errorf("RegisterSink() error = %v", err)
			}

			for _, ev := range tt.events {
				evs.Emit(ev)
			}

			select {
			case <-time.After(100 * time.Millisecond):
				t.Errorf("not all events recorded in time")
			case <-waitGroupDone(wg):
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, scenario(tt))
	}
}
