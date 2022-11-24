package sink_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	logging_mock "inetmock.icb4dc0.de/inetmock/internal/mock/logging"
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
	"inetmock.icb4dc0.de/inetmock/pkg/audit/sink"
	"inetmock.icb4dc0.de/inetmock/pkg/logging"
	"inetmock.icb4dc0.de/inetmock/pkg/wait"
)

func Test_logSink_OnSubscribe(t *testing.T) {
	t.Parallel()
	type fields struct {
		loggerSetup func(t *testing.T, wg *sync.WaitGroup) logging.Logger
	}
	type testCase struct {
		name   string
		fields fields
		events []*audit.Event
	}
	tests := []testCase{
		{
			name: "Get a single log line",
			fields: fields{
				loggerSetup: func(t *testing.T, wg *sync.WaitGroup) logging.Logger {
					t.Helper()
					ctrl := gomock.NewController(t)
					loggerMock := logging_mock.NewMockLogger(ctrl)

					loggerMock.
						EXPECT().
						With(gomock.Any()).
						Return(loggerMock)

					loggerMock.
						EXPECT().
						Info("handled request", gomock.Any()).
						Do(func(_ string, _ ...zap.Field) {
							wg.Done()
						}).
						Times(1)

					return loggerMock
				},
			},
			events: testEvents()[:1],
		},
		{
			name: "Get multiple events",
			fields: fields{
				loggerSetup: func(t *testing.T, wg *sync.WaitGroup) logging.Logger {
					t.Helper()
					ctrl := gomock.NewController(t)
					loggerMock := logging_mock.NewMockLogger(ctrl)

					loggerMock.
						EXPECT().
						With(gomock.Any()).
						Return(loggerMock)

					loggerMock.
						EXPECT().
						Info("handled request", gomock.Any()).
						Do(func(_ string, _ ...zap.Field) {
							wg.Done()
						}).
						Times(2)

					return loggerMock
				},
			},
			events: testEvents(),
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				wg  sync.WaitGroup
				evs audit.EventStream
				err error
			)
			wg.Add(len(tt.events))

			if evs, err = audit.NewEventStream(logging.CreateTestLogger(t), audit.WithBufferSize(10)); err != nil {
				t.Errorf("NewEventStream() error = %v", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			if err = evs.RegisterSink(ctx, sink.NewLogSink(tt.fields.loggerSetup(t, &wg))); err != nil {
				t.Errorf("RegisterSink() error = %v", err)
			}

			for _, ev := range tt.events {
				evs.Emit(ev)
			}

			select {
			case <-time.After(100 * time.Millisecond):
				t.Errorf("not all events recorded in time")
			case <-wait.ForWaitGroupDone(&wg):
			}
		})
	}
}
