package audit

import (
	"encoding/binary"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"

	auditv1 "gitlab.com/inetmock/inetmock/pkg/audit/v1"
)

const (
	lengthBufferSize         = 4
	defaultPayloadBufferSize = 4096
)

type Reader interface {
	Read() (Event, error)
}

type EventReaderOption func(reader *eventReader)

func NewEventReader(source io.Reader, opts ...EventReaderOption) Reader {
	reader := &eventReader{
		source: source,
		lengthPool: &sync.Pool{
			New: func() any {
				buf := make([]byte, lengthBufferSize)
				return &buf
			},
		},
		payloadPool: &sync.Pool{
			New: func() any {
				buf := make([]byte, defaultPayloadBufferSize)
				return &buf
			},
		},
	}

	for _, opt := range opts {
		opt(reader)
	}

	return reader
}

type eventReader struct {
	lengthPool  *sync.Pool
	payloadPool *sync.Pool
	source      io.Reader
}

func (e eventReader) Read() (ev Event, err error) {
	lengthBufRef := e.lengthPool.Get().(*[]byte)
	defer e.lengthPool.Put(lengthBufRef)

	lengthBuf := *lengthBufRef
	if _, err = e.source.Read(lengthBuf); err != nil {
		return
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	var msgBuf []byte
	if length <= defaultPayloadBufferSize {
		bufRef := e.payloadPool.Get().(*[]byte)
		defer e.payloadPool.Put(bufRef)

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

	ev = NewEventFromProto(protoEv)

	return
}
