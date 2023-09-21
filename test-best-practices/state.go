package test_best_practices

import (
	"context"
	"errors"
	mapset "github.com/deckarep/golang-set"
)

type Event int64

const (
	deliveringEvent Event = iota + 1
	onHoldEvent
	reDeliveringEvent
	delivered
)

type State int64

const (
	initialState = iota
)

var (
	orderEventStateMap = map[Event]struct {
		currentStateSet mapset.Set
		targetState     State
	}{}
)

type Order struct {
	OrderId    string
	OrderState State
}

func GetOrder(orderId string) Order {
	return Order{}
}

func UpdateOrder(originalOrder, order Order) error {
	return nil
}

func UpdateOrderStateByEvent(ctx context.Context, orderId string, event Event) (err error) {
	order := GetOrder(orderId)
	stateMap, ok := orderEventStateMap[event]
	if !ok {
		return errors.New("event not exists")
	}

	if !stateMap.currentStateSet.Contains(order.OrderState) {
		return errors.New("current OrderState error")
	}

	updateOrder := Order{
		OrderId:    order.OrderId,
		OrderState: order.OrderState,
	}

	err = UpdateOrder(order, updateOrder)
	if err != nil {
		return err
	}
	return nil
}
