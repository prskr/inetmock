package netflow

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
)

func NewPerfEventReader(m *ebpf.Map, perCPUBufferSize int) (reader *PerfEventReader, err error) {
	reader = new(PerfEventReader)
	if reader.reader, err = perf.NewReader(m, perCPUBufferSize); err != nil {
		return nil, err
	} else {
		return reader, nil
	}
}

type PerfEventReader struct {
	reader *perf.Reader
}

func (r *PerfEventReader) Read() (pkt *Packet, err error) {
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

func (r *PerfEventReader) Close() error {
	return r.reader.Close()
}
