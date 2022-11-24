package netflow_test

import (
	"errors"
	"net"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"inetmock.icb4dc0.de/inetmock/netflow"
)

func TestPacketTransport_Start(t *testing.T) {
	t.Parallel()

	testErr := errors.New("there's something strange")
	tests := []struct {
		name    string
		elems   []any
		wantErr error
	}{
		{
			name: "Empty reader",
		},
		{
			name: "Single element reader",
			elems: []any{
				&netflow.Packet{
					SourceIP:   net.IPv4(1, 2, 3, 4),
					DestIP:     net.IPv4(10, 10, 1, 1),
					SourcePort: 21842,
					DestPort:   80,
					Transport:  netflow.ProtocolTCP,
				},
			},
		},
		{
			name: "Single error reader",
			elems: []any{
				testErr,
			},
			wantErr: testErr,
		},
		{
			name: "Packet and error reader",
			elems: []any{
				&netflow.Packet{
					SourceIP:   net.IPv4(1, 2, 3, 4),
					DestIP:     net.IPv4(10, 10, 1, 1),
					SourcePort: 21842,
					DestPort:   80,
					Transport:  netflow.ProtocolTCP,
				},
				testErr,
			},
			wantErr: testErr,
		},
		{
			name: "Multiple packets and error reader",
			elems: []any{
				&netflow.Packet{
					SourceIP:   net.IPv4(1, 2, 3, 4),
					DestIP:     net.IPv4(10, 10, 1, 1),
					SourcePort: 21842,
					DestPort:   80,
					Transport:  netflow.ProtocolTCP,
				},
				testErr,
				&netflow.Packet{
					SourceIP:   net.IPv4(9, 8, 7, 6),
					DestIP:     net.IPv4(10, 10, 1, 1),
					SourcePort: 13987,
					DestPort:   443,
					Transport:  netflow.ProtocolTCP,
				},
			},
			wantErr: testErr,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := packetReaderMockOf(tt.elems...)
			sink := new(packetSinkRecorder)
			transport := netflow.NewPacketTransport(reader, sink, netflow.ErrorSinkFunc(func(err error) {
				if !errors.Is(err, errMockEmpty) && !errors.Is(err, tt.wantErr) {
					t.Errorf("Error occurred during processing: %v", err)
				}
			}))
			go transport.Start()

			<-reader.Done()

			td.Cmp(t, sink.RecordedPackets(), td.Bag(filter(tt.elems, excludeError[any])...))
		})
	}
}

func filter[T any](elems []T, include func(t T) bool) []T {
	out := make([]T, 0, len(elems))
	for i := range elems {
		if elem := elems[i]; include(elem) {
			out = append(out, elems[i])
		}
	}

	return out
}

func excludeError[T any](t T) bool {
	_, isErr := any(t).(error)
	return !isErr
}
