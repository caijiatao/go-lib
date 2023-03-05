package delay_queue

import (
	"context"
	"log"
)

type consumer struct {
	queues []DelayQueue
}

// 消费的守护协程
// 1.多个队列每次消费到100个之后就跳到下一个队列，避免一个队列消息太多导致某些队列一直处于饥饿状态
func (c *consumer) start(ctx context.Context) error {
	go func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}

		err := c.run(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return nil
}

func (c *consumer) run(ctx context.Context) error {
	for {
		for _, queue := range c.queues {
			err := queue.Consumer(ctx)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	return nil
}
