//go:build linux

package consumers

import (
	"errors"
	"io"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"

	"inetmock.icb4dc0.de/inetmock/internal/pcap"
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

func (f *writerConsumer) Observe(pkg gopacket.Packet) {
	if f.packageWriter == nil {
		return
	}
	/*
	 * copy data and metadata
	 * this avoids the risk of manipulation before they are written
	 */
	var (
		buf = make([]byte, len(pkg.Data()))
		ci  = pkg.Metadata().CaptureInfo
	)
	copy(buf, pkg.Data())
	_ = f.packageWriter.WritePacket(ci, buf)
}

func (f *writerConsumer) Close() (err error) {
	if closer, ok := f.origWriter.(io.Closer); ok {
		err = errors.Join(err, closer.Close())
	}
	f.packageWriter = nil
	return
}
