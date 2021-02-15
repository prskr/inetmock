package audit

import (
	"encoding/binary"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
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
			New: func() interface{} {
				return make([]byte, lengthBufferSize)
			},
		},
		payloadPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, defaultPayloadBufferSize)
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
	lengthBuf := e.lengthPool.Get().([]byte)
	if _, err = e.source.Read(lengthBuf); err != nil {
		return
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	var msgBuf []byte
	if length <= defaultPayloadBufferSize {
		msgBuf = e.payloadPool.Get().([]byte)[:length]
	} else {
		msgBuf = make([]byte, length)
	}

	if _, err = e.source.Read(msgBuf); err != nil {
		return
	}

	protoEv := new(EventEntity)

	if err = proto.Unmarshal(msgBuf, protoEv); err != nil {
		return
	}

	ev = NewEventFromProto(protoEv)

	return
}
