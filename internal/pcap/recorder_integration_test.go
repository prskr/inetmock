//go:build sudo
// +build sudo

package pcap_test

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/maxatome/go-testdeep/td"
	"golang.org/x/net/context/ctxhttp"

	"inetmock.icb4dc0.de/inetmock/internal/pcap"
	"inetmock.icb4dc0.de/inetmock/internal/pcap/consumers"
	"inetmock.icb4dc0.de/inetmock/internal/test"
	"inetmock.icb4dc0.de/inetmock/internal/test/integration"
)

const (
	simulateTimeout  = 200 * time.Millisecond
	recordingTimeout = 500 * time.Millisecond
)

func Test_recorder_CompleteWorkflow(t *testing.T) {
	t.Parallel()

	srv := integration.NewTestHTTPServer(t, []string{
		`=> Status(204)`,
	}, nil)

	listener := test.NewTCPListener(t, "127.0.0.1:0")
	go srv.Listen(t, listener)
	client := test.HTTPClientForListener(t, listener)

	listenerPort := uint16(listener.Addr().(*net.TCPAddr).Port)

	recorder := pcap.NewRecorder()
	t.Cleanup(func() {
		td.CmpNoError(t, recorder.Close())
	})
	buffer := newSyncBuffer()
	var err error
	var inMemConsumer pcap.Consumer
	if inMemConsumer, err = consumers.NewWriterConsumer("InMem", buffer); err != nil {
		t.Errorf("NewWriterConsumer() error = %v", err)
		return
	}

	recordCtx, recordCancel := context.WithTimeout(context.Background(), recordingTimeout)
	t.Cleanup(recordCancel)

	var result *pcap.StartRecordingResult
	if result, err = recorder.StartRecording(recordCtx, "lo", inMemConsumer); err != nil {
		t.Errorf("StartRecording() error = %v", err)
		return
	}

	td.Cmp(t, result, &pcap.StartRecordingResult{ConsumerKey: "lo:InMem"})

	simulateCtx, simulateCancel := context.WithTimeout(context.Background(), simulateTimeout)
	t.Cleanup(simulateCancel)
	simulateTraffic(simulateCtx, t, client)
	<-buffer.Closed

	var pcapReader *pcapgo.Reader
	if pcapReader, err = pcapgo.NewReader(buffer); err != nil {
		t.Errorf("pcapgo.NewReader() error = %v", err)
		return
	}

	packetSource := gopacket.NewZeroCopyPacketSource(pcapReader, layers.LayerTypeEthernet, gopacket.WithLazy(true), gopacket.WithPool(true))

	for pkg := range packetSource.Packets(context.Background()) {
		ip4LayerRaw := pkg.Layer(layers.LayerTypeIPv4)
		tcpLayerRaw := pkg.Layer(layers.LayerTypeTCP)
		if tcpLayerRaw == nil || ip4LayerRaw == nil {
			continue
		}

		ip4Layer, _ := ip4LayerRaw.(*layers.IPv4)
		tcpLayer, _ := tcpLayerRaw.(*layers.TCP)

		if ip4Layer.DstIP.IsLoopback() && uint16(tcpLayer.DstPort) == listenerPort {
			t.Logf("found one of the sample requests: %s", pkg.String())
			return
		}
	}
	t.Errorf("Couldn't find any of the sample requets")
}

func simulateTraffic(ctx context.Context, tb testing.TB, client *http.Client) {
	tb.Helper()
	for ctx.Err() == nil {
		_, err := ctxhttp.Get(ctx, client, "http://gitlab.com/")
		if errors.Is(err, context.DeadlineExceeded) {
			return
		}
		td.CmpNoError(tb, err)
	}
}

func newSyncBuffer() *syncBuffer {
	return &syncBuffer{
		Closed: make(chan bool),
		lock:   new(sync.Mutex),
		buf:    bytes.NewBuffer(nil),
	}
}

type syncBuffer struct {
	Closed chan bool
	lock   sync.Locker
	buf    *bytes.Buffer
}

func (s *syncBuffer) Close() error {
	select {
	case s.Closed <- true:
	case <-time.After(10 * time.Millisecond):
	}
	return nil
}

func (s *syncBuffer) Write(p []byte) (n int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) Read(p []byte) (n int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.buf.Read(p)
}
