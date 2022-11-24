package netflow

import (
	"errors"

	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
)

func NewPacketTransport(reader PacketReader, packetSink PacketSink, errorSink ErrorSink) *PacketTransport {
	return &PacketTransport{
		reader:     reader,
		PacketSink: packetSink,
		ErrorSink:  errorSink,
	}
}

type PacketTransport struct {
	PacketSink
	ErrorSink
	done   chan struct{}
	reader PacketReader
}

func (t *PacketTransport) Start() {
	t.done = make(chan struct{})

	for {
		if pkt, err := t.reader.Read(); err != nil {
			if errors.Is(err, ringbuf.ErrClosed) || errors.Is(err, perf.ErrClosed) {
				return
			}
			if errSink := t.ErrorSink; errSink != nil {
				errSink.OnError(err)
			}
			continue
		} else {
			t.PacketSink.OnObservedPacket(pkt)
			pkt.Dispose()
		}

		select {
		case _, more := <-t.done:
			if !more {
				return
			}
		default:
			continue
		}
	}
}

func (t *PacketTransport) Close() error {
	if t.done != nil {
		close(t.done)
	}

	return t.reader.Close()
}
