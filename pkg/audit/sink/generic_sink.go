package sink

import (
	"inetmock.icb4dc0.de/inetmock/pkg/audit"
)

func NewNoOpSink(name string) audit.Sink {
	return NewGenericSink(name, func(_ *audit.Event) {})
}

func NewGenericSink(name string, consumer func(ev *audit.Event)) audit.Sink {
	return &genericSink{
		name:     name,
		consumer: consumer,
	}
}

type genericSink struct {
	name     string
	consumer func(ev *audit.Event)
}

func (g genericSink) Name() string {
	return g.name
}

func (g genericSink) OnEvent(ev *audit.Event) {
	g.consumer(ev)
}
