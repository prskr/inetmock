package netflow_test

import (
	"bytes"
	_ "embed"
	"errors"
	"sync"
	"testing"

	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"
	"golang.org/x/exp/slices"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

var (
	_            netflow.PacketReader = (*packetReaderMock)(nil)
	_            netflow.PacketSink   = (*packetSinkRecorder)(nil)
	errMockEmpty                      = errors.New("mock is empty")

	//go:embed ebpf/tests.o
	testsEBPFProgram []byte
)

func packetReaderMockOf(elems ...any) *packetReaderMock {
	mock := &packetReaderMock{
		done:    make(chan struct{}),
		backlog: make(chan any, len(elems)),
	}

	for i := range elems {
		mock.backlog <- elems[i]
	}

	return mock
}

type packetReaderMock struct {
	backlog chan any
	done    chan struct{}
}

func (p *packetReaderMock) Read() (*netflow.Packet, error) {
	select {
	case e, more := <-p.backlog:
		if !more {
			return nil, errMockEmpty
		}
		switch u := e.(type) {
		case netflow.Packet:
			return &u, nil
		case *netflow.Packet:
			return u, nil
		case error:
			return nil, u
		}
	default:
		_ = p.Close()
	}
	return nil, errMockEmpty
}

func (p *packetReaderMock) Done() <-chan struct{} {
	return p.done
}

func (p *packetReaderMock) Close() error {
	close(p.backlog)
	close(p.done)
	return nil
}

type packetSinkRecorder struct {
	lock     sync.Mutex
	recorded []*netflow.Packet
}

func (p *packetSinkRecorder) RecordedPackets() []*netflow.Packet {
	p.lock.Lock()
	defer p.lock.Unlock()

	return slices.Clone(p.recorded)
}

func (p *packetSinkRecorder) OnObservedPacket(pkt *netflow.Packet) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.recorded = append(p.recorded, pkt)
}

func PrepareTestManager(tb testing.TB, nicName string, mode netflow.MonitorMode) *manager.Manager {
	tb.Helper()

	var (
		mgr  = new(manager.Manager)
		opts manager.Options
	)

	switch mode {
	case netflow.MonitorModePerfEvent:
		mgr.Probes = []*manager.Probe{
			{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFSection:  "classifier/perf-tests",
					EBPFFuncName: "emit_test_events_perf",
				},

				IfName:           nicName,
				NetworkDirection: manager.Egress,
			},
		}
		opts.ExcludedFunctions = append(opts.ExcludedFunctions, "emit_test_events_ring_buf")
	case netflow.MonitorModeRingBuf:
		mgr.Probes = []*manager.Probe{
			{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFSection:  "classifier/ringbuf-tests",
					EBPFFuncName: "emit_test_events_ring_buf",
				},

				IfName:           nicName,
				NetworkDirection: manager.Egress,
			},
		}
		opts.ExcludedFunctions = append(opts.ExcludedFunctions, "emit_test_events_perf")
		opts.MapSpecEditors = map[string]manager.MapSpecEditor{
			"observed_packets": {
				Type:       ebpf.RingBuf,
				MaxEntries: 1 << 24,
				EditorFlag: manager.EditType | manager.EditMaxEntries,
			},
		}
	case netflow.MonitorModeUnspecified:
		fallthrough
	default:
		tb.Fatalf("Monitor mode %v invalid", mode)
	}

	if err := mgr.InitWithOptions(bytes.NewReader(testsEBPFProgram), opts); err != nil {
		tb.Fatalf("Failed to initialize manager: %v", err)
	}

	if err := mgr.Start(); err != nil {
		tb.Fatalf("Failed to start manager: %v", err)
	}

	tb.Cleanup(func() {
		if err := mgr.Stop(manager.CleanAll); err != nil {
			tb.Errorf("Failed to stop manager: %v", err)
		}
	})

	return mgr
}

func ErrorSink(tb testing.TB) netflow.ErrorSinkOption {
	tb.Helper()
	return netflow.ErrorSinkOption{
		ErrorSink: netflow.ErrorSinkFunc(func(err error) {
			tb.Errorf("Error occurred during processing: %v", err)
		}),
	}
}
