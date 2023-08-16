package im

type Message struct {
	Key string
	Msg []byte
}

func NewMessage(key string, msg []byte) *Message {
	return &Message{
		Key: key,
		Msg: msg,
	}
}
