package goroutines

import (
	"testing"
	"time"
)

func TestBufChan(t *testing.T) {

	// 无buffer的channel, 接收一个值之后, 立刻被block
	ch := make(chan int)

	// 测试无buffer channel的行为
	noBufChanPlay(ch)

	// 这个与 make(chan int) 不同点在于, 它是有buffer的channel
	// 只有同时hold住2个值的时候, 投递者才会被block
	chBuf := make(chan int, 2)
	_ = chBuf // 暂时标记为已使用，避免编译错误
}

// noBufChanPlay -> 多次尝试可以发现, 这里取出消费先后次序是不能确保的
func noBufChanPlay(ch chan int) {
	// 异步投递一个
	go func(ch chan int) {
		ch <- 1
	}(ch)

	// 异步投递一个
	go func(ch chan int) {
		ch <- 3
	}(ch)

	// 等待1s
	time.Sleep(1 * time.Second)

	// 持续消费channel中的值，直到没有更多值可消费
	for {
		select {
		case value := <-ch:
			println("Consumed value:", value)
		default:
			// 如果channel中没有值了，就退出循环
			println("No more values in channel, exiting")
			return
		}
	}
}
