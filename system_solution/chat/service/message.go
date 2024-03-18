package service

import (
	"bytes"
	"compress/zlib"
	"io"
)

type Message struct {
	FromUser    int64  `json:"from_user"`
	ToUser      int64  `json:"to_user"`
	MessageBody string `json:"message_body"`
}

func NewPushMessageSuccessResp() []byte {
	return []byte(`{"status": "success"}`)
}

func NewPushMessageFailResp() []byte {
	return []byte(`{"status": "fail"}`)
}

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
