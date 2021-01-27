package sink

import (
	"gitlab.com/inetmock/inetmock/pkg/audit"
)

func NewNoOpSink(name string) audit.Sink {
	return NewGenericSink(name, func(_ audit.Event) {})
}

func NewGenericSink(name string, consumer func(ev audit.Event)) audit.Sink {
	return &genericSink{
		name:     name,
		consumer: consumer,
	}
}

type genericSink struct {
	name     string
	consumer func(ev audit.Event)
}

func (g genericSink) Name() string {
	return g.name
}

func (g genericSink) OnSubscribe(evs <-chan audit.Event) {
	go func(consumer func(ev audit.Event), evs <-chan audit.Event) {
		for ev := range evs {
			consumer(ev)
		}
	}(g.consumer, evs)
}
