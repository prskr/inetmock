//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/audit/writer.mock.go -package=audit_mock

package audit

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"google.golang.org/protobuf/proto"
)

var ErrValueMostNotBeNil = errors.New("event value must not be nil")

type Writer interface {
	io.Closer
	Write(ev *Event) error
}

type EventWriterOption func(writer *eventWriter)

func NewEventWriter(target io.Writer, opts ...EventWriterOption) Writer {
	writer := &eventWriter{
		target: target,
		lengthPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, lengthBufferSize)
			},
		},
	}

	for _, opt := range opts {
		opt(writer)
	}

	return writer
}

type eventWriter struct {
	lengthPool *sync.Pool
	target     io.Writer
}

type syncer interface {
	Sync() error
}

func (e eventWriter) Write(ev *Event) (err error) {
	if ev == nil {
		return ErrValueMostNotBeNil
	}
	var bytes []byte

	if bytes, err = proto.Marshal(ev.ProtoMessage()); err != nil {
		return
	}
	buf := e.lengthPool.Get().([]byte)
	binary.BigEndian.PutUint32(buf, uint32(len(bytes)))

	if _, err = e.target.Write(buf); err != nil {
		return
	}
	if _, err = e.target.Write(bytes); err != nil {
		return
	}
	if syncerInst, ok := e.target.(syncer); ok {
		err = syncerInst.Sync()
	}

	return
}

func (e eventWriter) Close() error {
	if closer, isCloser := e.target.(io.Closer); isCloser {
		return closer.Close()
	}
	return nil
}
