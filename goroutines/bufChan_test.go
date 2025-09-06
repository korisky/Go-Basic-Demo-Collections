package goroutines

import (
	"fmt"
	"testing"
	"time"
)

// TestNoBufChan 没有使用buf, 可以得到结果发现finish基本都慢于consume (因为chan只有1的capacity, 另一个投递就block住)
func TestNoBufChan(t *testing.T) {
	// 无buffer的channel, 接收一个值之后, 立刻被block
	ch := make(chan int)
	chanPlay(ch)
}

// TestBufChan 使用buf, 可以得到结果发现finish基本都早于consume（因为chan有buffer, 投递都没有block）
func TestBufChan(t *testing.T) {
	// 只有同时hold住2个值的时候, 投递者才会被block
	chBuf := make(chan int, 2)
	chanPlay(chBuf)
}

// chanPlay -> 多次尝试可以发现, 这里取出消费先后次序是不能确保的
func chanPlay(ch chan int) {

	// 异步投递一个
	go func(ch chan int) {
		ch <- 1
		fmt.Println("chan input 1 finish")
	}(ch)

	// 异步投递一个
	go func(ch chan int) {
		ch <- 3
		fmt.Println("chan input 3 finish")
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
