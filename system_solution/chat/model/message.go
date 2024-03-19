package model

import (
	"bytes"
	"compress/zlib"
	"io"
	"time"
)

type Message struct {
	Id         int64     `json:"id" gorm:"id"`
	FromUser   int64     `json:"from_user" gorm:"from_user"`
	ToUser     int64     `json:"to_user" gorm:"to_user"`
	Content    string    `json:"content" gorm:"content"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
}

type GroupMessage struct {
	Id         int64     `json:"id" gorm:"id"`
	GroupId    int64     `json:"group_id" gorm:"group_id"`
	FromUser   int64     `json:"from_user" gorm:"from_user"`
	Content    string    `json:"content" gorm:"content"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
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
