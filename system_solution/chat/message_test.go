package chat

import (
	"testing"
)

func TestCompressMessage(t *testing.T) {
	message := "hello, this is a message ,i am test the compress message, try the long message zip"
	compressed, err := CompressMessage(message)
	if err != nil {
		t.Error(err)
	}
	// size of the compressed
	if len(compressed) >= len([]byte(message)) {
		t.Errorf("Expected compressed message to be smaller than %d, got %d", len(message), len(compressed))
	}
	decompressed, err := DecompressMessage(compressed)
	if err != nil {
		t.Error(err)
	}
	if decompressed != message {
		t.Errorf("Expected %s, got %s", message, decompressed)
	}
}
