package sink

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/audit"
)

func NewGRPCSink(ctx context.Context, name string, consumer func(ev audit.Event)) audit.Sink {
	return &grpcSink{
		name:     name,
		ctx:      ctx,
		consumer: consumer,
	}
}

type grpcSink struct {
	name     string
	ctx      context.Context
	consumer func(ev audit.Event)
}

func (g grpcSink) Name() string {
	return g.name
}

func (g grpcSink) OnSubscribe(evs <-chan audit.Event, handle audit.CloseHandle) {
	go func(ctx context.Context, consumer func(ev audit.Event), evs <-chan audit.Event, handle audit.CloseHandle) {
		for {
			select {
			case ev := <-evs:
				consumer(ev)
			case <-ctx.Done():
				handle()
				return
			}
		}
	}(g.ctx, g.consumer, evs, handle)
}
