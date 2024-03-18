package test_best_practices

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderStateByEvent(t *testing.T) {
	type args struct {
		ctx     context.Context
		orderId string
		event   Event
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		initStubs func() (reset func())
	}{
		{
			name: "",
			args: args{
				ctx:     context.Background(),
				orderId: "orderId1",
				event:   onHoldEvent,
			},
			wantErr: nil,
			initStubs: func() (reset func()) {
				patches := gomonkey.ApplyFunc(GetOrder, func(orderId string) Order {
					return Order{
						OrderId:    orderId,
						OrderState: delivering,
					}
				})
				return func() {
					patches.Reset()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. 把需要mock 的方法mock掉
			reset := tt.initStubs()
			defer reset()
			// 2. 调用需要测试的方法
			err := UpdateOrderStateByEvent(tt.args.ctx, tt.args.orderId, tt.args.event)
			assert.Nil(t, err)
		})
	}
}
