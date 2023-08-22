package base

type Encoder interface {
	Encode(v any) ([]byte, error)
}

type Decoder interface {
	Decode(data []byte, v any) error
}

type Coder interface {
	Encoder
	Decoder
}
