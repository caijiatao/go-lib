package good_t

import "errors"

const (
	Pending = iota
	Paid
	Cancel
	Shipped
	Complete
)

type State int

type Order struct {
	Id         string
	State      State
	OrderState OrderState
}

type OrderAction int

const (
	OrderActionPay OrderAction = iota + 1
	OrderActionShip
	OrderActionComplete
	OrderActionCancel
)

var (
	orderAction2ValidState = map[OrderAction][]State{
		OrderActionPay:      {Pending, Shipped},
		OrderActionShip:     {Paid},
		OrderActionComplete: {Shipped, Paid},
		OrderActionCancel:   {Pending},
	}
)

func ValidOrderState(order Order, action OrderAction) bool {
	validStates, ok := orderAction2ValidState[action]
	if !ok {
		return false
	}

	for _, state := range validStates {
		if order.State == state {
			return true
		}
	}

	return false
}

func CreateOrder() Order {
	// 在数据库中创建订单

	return Order{
		State: Pending,
	}
}

func PaidOrder(order Order) error {
	// 未支付的订单才能支付
	if ValidOrderState(order, OrderActionPay) {
		return errors.New("order state is not pending")
	}

	// 更新订单状态到数据库
	order.State = Paid
	// ...

	return nil
}

func CancelOrder(order Order) error {
	// 只有订单状态为Pending时才能取消
	if ValidOrderState(order, OrderActionCancel) {
		return errors.New("order state is not pending")
	}

	// 更新订单状态到数据库
	order.State = Cancel
	// ...

	return nil
}

func ShippedOrder(order Order) error {
	// 只有订单状态为Paid时才能发货
	if ValidOrderState(order, OrderActionShip) {
		return errors.New("order state is not paid")
	}

	// 更新订单状态到数据库
	order.State = Shipped
	// ...

	return nil
}

func CompleteOrder(order Order) error {
	// 只有订单状态为Shipped时才能完成
	if ValidOrderState(order, OrderActionComplete) {
		return errors.New("order state is not shipped")
	}

	// 更新订单状态到数据库
	order.State = Complete
	// ...

	return nil
}
