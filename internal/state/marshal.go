package state

type Encoder interface {
	Encode(v interface{}) (data []byte, err error)
}

type Decoder interface {
	Decode(data []byte, v interface{}) error
}

type EncoderDecoder interface {
	Encoder
	Decoder
}
