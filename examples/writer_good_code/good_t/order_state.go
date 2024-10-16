package good_t

import (
	"errors"
	"fmt"
)

// OrderState 定义订单状态行为
type OrderState interface {
	Paid(order *Order) error
	Cancel(order *Order) error
	GetState() int
}

// PendingState 待支付状态
type PendingState struct{}

func (s *PendingState) Paid(order *Order) error {
	// 更新订单状态到数据库
	order.OrderState = &PaidState{}
	fmt.Println("订单已支付")
	return nil
}

func (s *PendingState) Cancel(order *Order) error {
	// 更新订单状态到数据库
	order.OrderState = &CanceledState{}
	fmt.Println("订单已取消")
	return nil
}

func (s *PendingState) GetState() int {
	return Pending
}

// PaidState 已支付状态
type PaidState struct{}

func (s *PaidState) Paid(order *Order) error {
	return errors.New("订单已支付，不能再次支付")
}

func (s *PaidState) Cancel(order *Order) error {
	return errors.New("订单已支付，不能取消")
}

func (s *PaidState) GetState() int {
	return Paid
}

// CanceledState 已取消状态
type CanceledState struct{}

func (s *CanceledState) Paid(order *Order) error {
	return errors.New("订单已取消，不能支付")
}

func (s *CanceledState) Cancel(order *Order) error {
	return errors.New("订单已取消，不能重复取消")
}

func (s *CanceledState) GetState() int {
	return Cancel
}

// StateCreateOrder 创建一个新订单，默认状态为 Pending
func StateCreateOrder() *Order {
	return &Order{
		Id:         "123456",
		OrderState: &PendingState{},
	}
}

// StatePaidOrder 支付订单
func StatePaidOrder(order *Order) error {
	return order.OrderState.Paid(order)
}

// StateCancelOrder 取消订单
func StateCancelOrder(order *Order) error {
	return order.OrderState.Cancel(order)
}
