package goroutines

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultSelect(t *testing.T) {

	start := time.Now()
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)

	elapsed := func() time.Duration {
		return time.Since(start).Round(time.Millisecond)
	}

	for {
		select {
		// channel A
		case <-tick:
			fmt.Printf("[%6s] tick.\n", elapsed())
		// channel B
		case <-boom:
			fmt.Printf("[%6s] BOOM!\n", elapsed())
			return
		// 如果for循环执行的够快, 但是2个channel都还没有内容
		// 直接走入default
		default:
			fmt.Printf("[%6s]       .\n", elapsed())
			time.Sleep(50 * time.Millisecond)
		}
	}

}
