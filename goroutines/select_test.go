package goroutines

import (
	"fmt"
	"testing"
)

func TestFibonacciSelect(t *testing.T) {
	c := make(chan int)
	quit := make(chan int)

	// 异步进行数据存放
	go func() {
		for range 10 {
			fmt.Println(<-c)
		}
		// 这里是明确告知退出channel, 需要进行停止
		quit <- 0
	}()

	// for循环取值处理
	fibonacciSelect(c, quit)
}

func fibonacciSelect(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		// 正常从chanel取出值
		case c <- x:
			x, y = y, x+y
		// 这里是经典的使用方式, 当quit这个channel有值,
		// 则进行退出. 这里由go自己判断是否接受到了quit
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
