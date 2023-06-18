package coder

type Coder interface {
	Encode(msg IMessage, body []byte) (buffer []byte, err error)
	Decode(msg IMessage, buffer []byte) (body []byte, err error)
}
