package goroutines

import (
	"fmt"
	"testing"
)

func TestChanStuff(t *testing.T) {

	// init the slice (array: [123]int, slice: []int)
	s := []int{7, 2, 8, -9, 4, 0, 5, 12, 3354, 1234, 52, -2523, 123}

	// 这里的channel是没有buffer的, 意味着一旦接收一个值之后
	// 一定要被取出消费, 之后才能接收另一个值
	// 也就是位置被占了, 就会进入block状态, 投递者的block
	c := make(chan int)

	// async summing
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	// retrieve from chan
	x, y := <-c, <-c

	fmt.Println(x, y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send result to channel
}
