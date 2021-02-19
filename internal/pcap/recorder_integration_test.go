package pcap_test

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"

	"gitlab.com/inetmock/inetmock/internal/pcap"
	"gitlab.com/inetmock/inetmock/internal/pcap/consumers"
)

const (
	simulateTimeout = 30 * time.Second
)

func Test_recorder_CompleteWorkflow(t *testing.T) {
	var recorder = pcap.NewRecorder()
	var buffer = bytes.NewBuffer(nil)
	var err error
	var inMemConsumer pcap.Consumer
	if inMemConsumer, err = consumers.NewWriterConsumer("InMem", buffer); err != nil {
		t.Errorf("NewWriterConsumer() error = %v", err)
		return
	}

	recordCtx, recordCancel := context.WithCancel(context.Background())
	defer recordCancel()
	if err = recorder.StartRecording(recordCtx, "lo", inMemConsumer); err != nil {
		t.Errorf("StartRecording() error = %v", err)
		recordCancel()
		return
	}

	simulateCtx, simulateCancel := context.WithTimeout(context.Background(), simulateTimeout)
	defer simulateCancel()
	simulateTraffic(simulateCtx, 100)

	recordCancel()

	var pcapReader *pcapgo.Reader
	if pcapReader, err = pcapgo.NewReader(buffer); err != nil {
		t.Errorf("pcapgo.NewReader() error = %v", err)
		return
	}

	packetSource := gopacket.NewPacketSource(pcapReader, layers.LayerTypeEthernet)
	packetSource.NoCopy = true
	packetSource.Lazy = true

	for pkg := range packetSource.Packets() {
		ip4LayerRaw := pkg.Layer(layers.LayerTypeIPv4)
		tcpLayerRaw := pkg.Layer(layers.LayerTypeTCP)
		if tcpLayerRaw == nil || ip4LayerRaw == nil {
			continue
		}

		ip4Layer, _ := ip4LayerRaw.(*layers.IPv4)
		tcpLayer, _ := tcpLayerRaw.(*layers.TCP)

		if ip4Layer.DstIP.IsLoopback() && tcpLayer.DstPort == 8181 {
			t.Logf("found one of the sample requests: %s", pkg.String())
			return
		}
	}
	t.Errorf("Couldn't find any of the sample requets")
}

func simulateTraffic(ctx context.Context, numberOfRequests int) {
	for i := 0; i < numberOfRequests; i++ {
		sampleURL, _ := url.Parse("http://127.0.0.1:8181")
		req := (&http.Request{
			Method: http.MethodGet,
			URL:    sampleURL,
		}).WithContext(ctx)
		_, _ = http.DefaultClient.Do(req)
	}
}
