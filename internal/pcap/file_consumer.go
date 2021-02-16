package pcap

import (
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"go.uber.org/multierr"
)

type writerConsumer struct {
	name           string
	origWriter     io.Writer
	packageWriter  *pcapgo.Writer
	snapshotLength int32
	linkType       layers.LinkType
}

func (f *writerConsumer) Init(params CaptureParameters) {
	f.snapshotLength = params.SnapshotLength
	f.linkType = params.LinkType
	_ = f.packageWriter.WriteFileHeader(uint32(params.SnapshotLength), params.LinkType)
}

func NewWriterConsumer(name string, writer io.Writer) (consumer Consumer, err error) {
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
	if pkg.ApplicationLayer() != nil {
		_ = f.packageWriter.WritePacket(pkg.Metadata().CaptureInfo, pkg.Data())
	}
}

func (f *writerConsumer) Close() (err error) {
	if closer, ok := f.origWriter.(io.Closer); ok {
		err = multierr.Append(err, closer.Close())
	}
	f.packageWriter = nil
	return nil
}
