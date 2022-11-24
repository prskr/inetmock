package audit

import (
	"encoding/binary"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"

	auditv1 "inetmock.icb4dc0.de/inetmock/pkg/audit/v1"
)

const (
	lengthBufferSize         = 4
	defaultPayloadBufferSize = 4096
)

var (
	lengthPool = sync.Pool{
		New: func() any {
			buf := make([]byte, lengthBufferSize)
			return &buf
		},
	}

	payloadPool = sync.Pool{
		New: func() any {
			buf := make([]byte, defaultPayloadBufferSize)
			return &buf
		},
	}
)

type Reader interface {
	Read() (*Event, error)
}

type EventReaderOption func(reader *eventReader)

func NewEventReader(source io.Reader, opts ...EventReaderOption) Reader {
	reader := &eventReader{
		source: source,
	}

	for _, opt := range opts {
		opt(reader)
	}

	return reader
}

type eventReader struct {
	source io.Reader
}

func (e eventReader) Read() (ev *Event, err error) {
	lengthBufRef := lengthPool.Get().(*[]byte)

	//nolint:staticcheck // make sure slice size is still buffer size
	defer func() {
		b := *lengthBufRef
		b = b[:lengthBufferSize]
		lengthPool.Put(lengthBufRef)
	}()

	lengthBuf := *lengthBufRef
	if _, err = e.source.Read(lengthBuf); err != nil {
		return
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	var msgBuf []byte
	if length <= defaultPayloadBufferSize {
		bufRef := payloadPool.Get().(*[]byte)

		//nolint:staticcheck // make sure slice size is still buffer size
		defer func() {
			b := *bufRef
			b = b[:defaultPayloadBufferSize]
			payloadPool.Put(bufRef)
		}()

		msgBuf = (*bufRef)[:length]
	} else {
		msgBuf = make([]byte, length)
	}

	if _, err = e.source.Read(msgBuf); err != nil {
		return
	}

	protoEv := new(auditv1.EventEntity)

	if err = proto.Unmarshal(msgBuf, protoEv); err != nil {
		return
	}

	return NewEventFromProto(protoEv), nil
}
