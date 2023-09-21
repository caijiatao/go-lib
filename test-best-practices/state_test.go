package test_best_practices

import (
	"context"
	"github.com/agiledragon/gomonkey/v2"
	"testing"
)

func TestOrderStateByEvent(t *testing.T) {
	type args struct {
		ctx     context.Context
		orderId string
		event   Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				OrderId:    "test-order",
				OrderState: 0,
			}
			gomonkey.ApplyFunc(GetOrder, func(orderId string) Order {
				order.OrderId = orderId
				return order
			})

			gomonkey.ApplyFunc(UpdateOrder, func() {})

		})
	}
}
