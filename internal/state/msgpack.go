package state

import "github.com/vmihailenco/msgpack/v5"

type MsgPackEncoding struct{}

func (m MsgPackEncoding) Encode(v any) (data []byte, err error) {
	return msgpack.Marshal(v)
}

func (m MsgPackEncoding) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
