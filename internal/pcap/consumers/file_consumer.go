// +build linux
//go:build linux

package consumers

import (
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"go.uber.org/multierr"

	"gitlab.com/inetmock/inetmock/internal/pcap"
)

const defaultSnapshotLength = 65536

var _ pcap.Consumer = (*writerConsumer)(nil)

type writerConsumer struct {
	name          string
	origWriter    io.Writer
	packageWriter *pcapgo.Writer
}

func (f *writerConsumer) Init() error {
	return f.packageWriter.WriteFileHeader(defaultSnapshotLength, layers.LinkTypeEthernet)
}

func NewWriterConsumer(name string, writer io.Writer) (consumer pcap.Consumer, err error) {
	consumer = &writerConsumer{
		name:          name,
		origWriter:    writer,
		packageWriter: pcapgo.NewWriter(writer),
	}
	return
}

func (f writerConsumer) Name() string {
	return f.name
}

func (f writerConsumer) Observe(pkg gopacket.Packet) {
	if f.packageWriter == nil {
		return
	}
	_ = f.packageWriter.WritePacket(pkg.Metadata().CaptureInfo, pkg.Data())
}

func (f *writerConsumer) Close() (err error) {
	if closer, ok := f.origWriter.(io.Closer); ok {
		multierr.AppendInvoke(&err, multierr.Close(closer))
	}
	f.packageWriter = nil
	return
}
