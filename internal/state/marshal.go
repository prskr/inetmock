package state

type Encoder interface {
	Encode(v any) (data []byte, err error)
}

type Decoder interface {
	Decode(data []byte, v any) error
}

type EncoderDecoder interface {
	Encoder
	Decoder
}
