package client

import "context"

type Client interface {
	Call(ctx context.Context, req interface{}, resp interface{}, opts ...Option) error
}

type MockClient struct{}

func (m *MockClient) Call(ctx context.Context, req interface{}, resp interface{}, opts ...Option) error {
	//TODO implement me
	panic("implement me")
}
