//go:build sudo

package netflow_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cilium/ebpf"
	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/netflow"
	"inetmock.icb4dc0.de/inetmock/pkg/wait"
)

//nolint:paralleltest // cannot test parallel with same interface
func TestPerfEventReader_Read(t *testing.T) {
	RemoveMemlock(t)
	go MakeSomeNoise(t, 10*time.Millisecond)

	tests := []struct {
		name          string
		wantPkt       td.TestDeep
		expectPackets int
		timeout       time.Duration
		mode          netflow.MonitorMode
		readerSetup   func(tb testing.TB, m *ebpf.Map) netflow.PacketReader
	}{
		{
			name: "Get a single packet - perf events",
			wantPkt: td.SuperBagOf(td.Struct(new(netflow.Packet), td.StructFields{
				"SourcePort": uint16(23876),
				"DestPort":   td.Any(uint16(53), uint16(80), uint16(443), uint16(853), uint16(3128), uint16(8080)),
			})),
			expectPackets: 1,
			mode:          netflow.MonitorModePerfEvent,
			timeout:       50 * time.Millisecond,
			readerSetup:   preparePerfEventReader,
		},
		{
			name: "Get multiple packets - perf events",
			wantPkt: td.SuperBagOf(td.Struct(new(netflow.Packet), td.StructFields{
				"SourcePort": uint16(23876),
				"DestPort":   td.Any(uint16(53), uint16(80), uint16(443), uint16(853), uint16(3128), uint16(8080)),
			})),
			expectPackets: 5,
			mode:          netflow.MonitorModePerfEvent,
			timeout:       100 * time.Millisecond,
			readerSetup:   preparePerfEventReader,
		},
		{
			name: "Get a single packet - ring buf",
			wantPkt: td.SuperBagOf(td.Struct(new(netflow.Packet), td.StructFields{
				"SourcePort": uint16(23876),
				"DestPort":   td.Any(uint16(53), uint16(80), uint16(443), uint16(853), uint16(3128), uint16(8080)),
			})),
			expectPackets: 1,
			mode:          netflow.MonitorModeRingBuf,
			timeout:       50 * time.Millisecond,
			readerSetup:   prepareRingBufReader,
		},
		{
			name: "Get multiple packets - ring buf",
			wantPkt: td.SuperBagOf(td.Struct(new(netflow.Packet), td.StructFields{
				"SourcePort": uint16(23876),
				"DestPort":   td.Any(uint16(53), uint16(80), uint16(443), uint16(853), uint16(3128), uint16(8080)),
			})),
			expectPackets: 5,
			mode:          netflow.MonitorModeRingBuf,
			timeout:       100 * time.Millisecond,
			readerSetup:   prepareRingBufReader,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mgr := PrepareTestManager(t, "lo", tt.mode)
			m, present, err := mgr.GetMap("observed_packets")
			if err != nil || !present {
				t.Fatalf("Failed to get observed_packets map")
			}

			var (
				wg         sync.WaitGroup
				gotPackets = make([]*netflow.Packet, 0, tt.expectPackets)
				reader     = tt.readerSetup(t, m)
			)

			wg.Add(tt.expectPackets)

			go func(wg *sync.WaitGroup) {
				for i := 0; i < tt.expectPackets; i++ {
					if pkt, err := reader.Read(); err != nil {
						t.Errorf("Read() error = %v", err)
						return
					} else {
						gotPackets = append(gotPackets, pkt)
					}
					wg.Done()
				}
			}(&wg)

			select {
			case <-wait.ForWaitGroupDone(&wg):
				td.Cmp(t, gotPackets, tt.wantPkt)
			case <-time.After(tt.timeout):
				t.Errorf("Did not complete in time")
			}
		})
	}
}

func preparePerfEventReader(tb testing.TB, m *ebpf.Map) netflow.PacketReader {
	tb.Helper()
	reader, err := netflow.NewPerfEventReader(m, 8)
	if err != nil {
		tb.Fatalf("Failed to init reader: %v", err)
	}

	return reader
}

func prepareRingBufReader(tb testing.TB, m *ebpf.Map) netflow.PacketReader {
	tb.Helper()
	reader, err := netflow.NewRingBufReader(m)
	if err != nil {
		tb.Fatalf("Failed to init reader: %v", err)
	}

	return reader
}
