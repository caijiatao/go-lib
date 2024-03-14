package model

import "time"

type Order struct {
	OrderId     string    `json:"orderId"`
	UserId      string    `json:"userId"`
	Price       string    `json:"price"`
	CreateTime  time.Time `json:"createTime"`
	PaymentTime time.Time
	User        User `json:"user"`
}
