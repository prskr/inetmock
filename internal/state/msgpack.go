package state

import "github.com/vmihailenco/msgpack/v5"

type MsgPackEncoding struct{}

func (m MsgPackEncoding) Encode(v interface{}) (data []byte, err error) {
	return msgpack.Marshal(v)
}

func (m MsgPackEncoding) Decode(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
