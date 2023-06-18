package coder

type IMessage interface {
	// SerializationType returns serialization type.
	SerializationType() int
}
