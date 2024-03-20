package model

import (
	"bytes"
	"compress/zlib"
	"io"
	"time"
)

type MessageStatus int

const (
	UnSend MessageStatus = iota + 1
	Send
	Received
	Read
)

// Message
// @Description: 用户消息，如果已经是已读的状态，则可以按需求保留一定天数，比如产品形态上需要同步近期消息，近期的时间如果是7天，则7天之前的消息已读的消息则可以归档
type Message struct {
	Id         int64         `json:"id" gorm:"id"`
	FromUser   int64         `json:"from_user" gorm:"from_user"`
	ToUser     int64         `json:"to_user" gorm:"to_user"`
	Content    string        `json:"content" gorm:"content"`
	Status     MessageStatus `json:"status" gorm:"status"`
	CreateTime time.Time     `json:"create_time" gorm:"create_time"`
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

type DeviceOptType int

const (
	DeviceOptNoPush DeviceOptType = iota + 1
	DeviceOptPush
)

type DeviceStatusType int

const (
	DeviceOffline DeviceStatusType = iota + 1
	DeviceOnline
	DeviceBusy
	DeviceLeave
	DeviceHide
)

// UserDeviceMessage
// @Description: 用户设备消息
type UserDeviceMessage struct {
	UserId          int64            `json:"user_id" gorm:"user_id"`
	DeviceId        string           `json:"device_id" gorm:"device_id"`
	CurMaxMessageId string           `json:"cur_max_message_id" gorm:"cur_max_message_id"`
	DeviceOpt       DeviceOptType    `json:"device_opt" gorm:"device_opt"`
	DeviceStatus    DeviceStatusType `json:"device_status" gorm:"device_status"`
	CreateTime      time.Time        `json:"create_time" gor:"create_time"`
}
