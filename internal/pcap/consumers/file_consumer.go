// +build linux

package consumers

import (
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcapgo"
	"gitlab.com/inetmock/inetmock/internal/pcap"
	"go.uber.org/multierr"
)

type writerConsumer struct {
	name          string
	origWriter    io.Writer
	packageWriter *pcapgo.Writer
}

func (f *writerConsumer) Init(params pcap.CaptureParameters) {
	_ = f.packageWriter.WriteFileHeader(65536, params.LinkType)
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
		err = multierr.Append(err, closer.Close())
	}
	f.packageWriter = nil
	return nil
}
