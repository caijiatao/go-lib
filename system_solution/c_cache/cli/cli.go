package cli

import "context"

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

type Result struct{}

func (r *Result) String() string {
	return "result 格式化结果"
}

func (c *Client) Do(ctx context.Context, cmd string, args ...string) (*Result, error) {
	return nil, nil
}

func (c *Client) Get(ctx context.Context, args ...string) (*Result, error) {
	return nil, nil
}

func (c *Client) Put(ctx context.Context, args ...string) (*Result, error) {
	return nil, nil
}

func (c *Client) Delete(ctx context.Context, args ...string) (*Result, error) {
	return nil, nil
}

func (c *Client) Ping(ctx context.Context, args ...interface{}) (*Result, error) {
	return nil, nil
}
