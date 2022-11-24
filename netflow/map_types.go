package netflow

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/valyala/bytebufferpool"
)

var _ encoding.BinaryMarshaler = (*valueCollection[int])(nil)

type valueCollection[T any] []T

func (col *valueCollection[T]) Set(idx int, val T) {
	(*col)[idx] = val
}

func (col valueCollection[T]) Get(idx int) T {
	return col[idx]
}

func (col valueCollection[T]) MarshalBinary() (data []byte, err error) {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if len(col) < 1 {
		return nil, nil
	}

	for i := range col {
		val := any(col[i])
		if m, ok := val.(encoding.BinaryMarshaler); ok {
			if d, err := m.MarshalBinary(); err != nil {
				return nil, err
			} else if _, err = buf.Write(d); err != nil {
				return nil, err
			}
		} else if err := binary.Write(buf, binary.LittleEndian, col[i]); err != nil {
			return nil, err
		}
	}

	tmp := buf.Bytes()
	data = make([]byte, len(tmp))
	copy(data, tmp)
	buf.Reset()

	return data, nil
}

func (col valueCollection[T]) UnmarshalBinary(data []byte) (err error) {
	if len(col) < 1 {
		return errors.New("cannot unmarshal into empty collection")
	}
	var binarySize int
	if m, ok := any(&col[0]).(BinaryCollectionUnmarshaler); ok {
		binarySize = m.BinarySize()
		if l := len(data); l%binarySize != 0 {
			return fmt.Errorf("data length %d is not a multiple of expected binary size %d", l, binarySize)
		}
	}

	if binarySize == 0 {
		return binary.Read(bytes.NewReader(data), binary.LittleEndian, col)
	}

	countElements := len(data) / binarySize

	if colLength := len(col); colLength < countElements {
		countElements = colLength
	}

	for i := 0; i < countElements; i++ {
		e := any(&col[i])
		if um, ok := e.(encoding.BinaryUnmarshaler); ok {
			if err := um.UnmarshalBinary(data[i*binarySize : (i+1)*binarySize]); err != nil {
				return err
			}
		} else {
			return errors.New("element does not implement encoding.BinaryUnmarshaler")
		}
	}

	return nil
}
