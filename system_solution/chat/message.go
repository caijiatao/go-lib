package chat

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressMessage(message string) ([]byte, error) {
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)

	_, err := writer.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

func DecompressMessage(compressed []byte) (string, error) {
	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decompressed), nil
}
