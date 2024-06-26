// Code generated by goctl. DO NOT EDIT.
// Source: send_svr.proto

package sender

import (
	"context"

	"chat-app-svr/rpc/send_svr/send"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	SendGroupMessageReply   = send.SendGroupMessageReply
	SendGroupMessageRequest = send.SendGroupMessageRequest
	SendMessageReply        = send.SendMessageReply
	SendMessageRequest      = send.SendMessageRequest

	Sender interface {
		SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error)
		SendGroupMessage(ctx context.Context, in *SendGroupMessageRequest, opts ...grpc.CallOption) (*SendGroupMessageReply, error)
	}

	defaultSender struct {
		cli zrpc.Client
	}
)

func NewSender(cli zrpc.Client) Sender {
	return &defaultSender{
		cli: cli,
	}
}

func (m *defaultSender) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error) {
	client := send.NewSenderClient(m.cli.Conn())
	return client.SendMessage(ctx, in, opts...)
}

func (m *defaultSender) SendGroupMessage(ctx context.Context, in *SendGroupMessageRequest, opts ...grpc.CallOption) (*SendGroupMessageReply, error) {
	client := send.NewSenderClient(m.cli.Conn())
	return client.SendGroupMessage(ctx, in, opts...)
}
