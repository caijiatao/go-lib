package delay_queue

import "context"

type DelayQueue interface {
	// Consumer
	// @Description:
	Consumer(ctx context.Context) error

	// Send
	// @Description: 延迟队列中发送消息
	Send(ctx context.Context) error
}
