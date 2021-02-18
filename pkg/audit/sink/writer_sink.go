package sink

import "gitlab.com/inetmock/inetmock/pkg/audit"

type WriterSinkOption func(sink *writerCloserSink)

var (
	WithCloseOnExit WriterSinkOption = func(sink *writerCloserSink) {
		sink.closeOnExit = true
	}
)

func NewWriterSink(name string, target audit.Writer, opts ...WriterSinkOption) audit.Sink {
	sink := &writerCloserSink{
		name:   name,
		target: target,
	}

	for _, opt := range opts {
		opt(sink)
	}

	return sink
}

type writerCloserSink struct {
	name        string
	target      audit.Writer
	closeOnExit bool
}

func (f writerCloserSink) Name() string {
	return f.name
}

func (f writerCloserSink) OnSubscribe(evs <-chan audit.Event) {
	go func(target audit.Writer, closeOnExit bool, evs <-chan audit.Event) {
		for e := range evs {
			ev := e
			_ = target.Write(&ev)
		}
		if closeOnExit {
			_ = target.Close()
		}
	}(f.target, f.closeOnExit, evs)
}
