package netflow

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/ringbuf"
)

func NewRingBufReader(m *ebpf.Map) (reader *RingBufReader, err error) {
	reader = new(RingBufReader)
	if reader.reader, err = ringbuf.NewReader(m); err != nil {
		return nil, err
	} else {
		return reader, nil
	}
}

type RingBufReader struct {
	reader *ringbuf.Reader
}

func (r *RingBufReader) Read() (pkt *Packet, err error) {
	pkt = packetPool.Get().(*Packet)
	defer func() {
		if err != nil {
			packetPool.Put(pkt)
		}
	}()

	if rec, err := r.reader.Read(); err != nil {
		return nil, err
	} else if err := pkt.UnmarshalBinary(rec.RawSample); err != nil {
		return nil, err
	} else {
		return pkt, nil
	}
}

func (r *RingBufReader) Close() error {
	return r.reader.Close()
}
