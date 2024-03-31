package main

import (
	"fmt"
	"time"
)

// OrderState 使用int类型进行保存
type OrderState int

// 提前声明常量为内存连续区域
const (
	Init OrderState = iota + 1
	Pending
	Success
	Failure
)

// validTransition 是一个Map, 用来确保状态扭转的方向
var validTransition = map[OrderState][]OrderState{
	Init:    {Pending},
	Pending: {Success, Failure},
	Success: {},
	Failure: {},
}

// Order 是订单的实体类型
type Order struct {
	OrderNo    string
	CreateTime time.Time
	OrderState OrderState
}

// Transition 判断状态扭转是否合法
func (o *Order) Transition(newState OrderState) error {
	for _, validNxtState := range validTransition[o.OrderState] {
		if newState == validNxtState {
			fmt.Printf("Order %s transited from %s to %s\n", o.OrderNo, o.OrderState, newState)
			o.OrderState = newState
			return nil
		}
	}
	return fmt.Errorf("Invalid transition from %v to %v for order %s\n", o.OrderState, newState, o.OrderNo)
}

func main() {

	order := Order{
		OrderNo:    "2341fsdf1",
		CreateTime: time.Now(),
		OrderState: Init,
	}

	if err := order.Transition(Pending); err != nil {
		fmt.Println(err)
	}

	if err := order.Transition(Init); err != nil {
		fmt.Println(err)
	}

}
