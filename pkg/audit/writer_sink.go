package audit

type WriterSinkOption func(sink *writerCloserSink)

var (
	WithCloseOnExit WriterSinkOption = func(sink *writerCloserSink) {
		sink.closeOnExit = true
	}
)

func NewWriterSink(name string, target Writer, opts ...WriterSinkOption) Sink {
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
	target      Writer
	closeOnExit bool
}

type syncer interface {
	Sync() error
}

func (f writerCloserSink) Name() string {
	return f.name
}

func (f writerCloserSink) OnSubscribe(evs <-chan Event) {
	go func(target Writer, closeOnExit bool, evs <-chan Event) {
		for ev := range evs {
			_ = target.Write(&ev)
		}
		if closeOnExit {
			_ = target.Close()
		}
	}(f.target, f.closeOnExit, evs)
}
