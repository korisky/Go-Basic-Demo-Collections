package goroutines

import (
	"fmt"
	"testing"
)

func TestCloseChan(t *testing.T) {
	// 创建buffer为10的channel
	c := make(chan int, 10)
	// 多一条协程, 进行叠加操作
	go fibonacci(cap(c), c)
	// 注意
	for i := range c {
		fmt.Println(i)
	}
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	// 这里主动关闭了channel, 而且需要注意只有sender可以close掉channel
	// 另外需要注意，外部配合了range可以很好的处理close的channel
	// 如果我们注释close，那么可以发现外部主线程会stuck，然后抛出deadlock
	close(c)
}
