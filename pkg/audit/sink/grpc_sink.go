package sink

import (
	"context"

	"gitlab.com/inetmock/inetmock/pkg/audit"
)

func NewGRPCSink(name string, consumer func(ev audit.Event)) audit.Sink {
	return &grpcSink{
		name:     name,
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

func (g grpcSink) OnSubscribe(evs <-chan audit.Event) {
	go func(consumer func(ev audit.Event), evs <-chan audit.Event) {
		for ev := range evs {
			consumer(ev)
		}
	}(g.consumer, evs)
}
