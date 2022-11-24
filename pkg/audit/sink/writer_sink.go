package sink

import "inetmock.icb4dc0.de/inetmock/pkg/audit"

type WriterSinkOption func(sink *writerCloserSink)

var WithCloseOnExit WriterSinkOption = func(sink *writerCloserSink) {
	sink.closeOnExit = true
}

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

func (f writerCloserSink) OnEvent(ev *audit.Event) {
	_ = f.target.Write(ev)
}

func (f writerCloserSink) Close() error {
	if f.closeOnExit {
		return f.target.Close()
	}
	return nil
}
