package goroutines

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// SafeCounter 不加锁的情况下, 会得到 ‘concurrent map writes’ 错误
type SafeCounter struct {
	mu sync.Mutex // 互斥锁
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key]++ // 使得线程同步
	// 可以使用c.mu.Unlock()在这里，进行解锁
	// 但推荐使用defer，使得所有操作更好的捆绑起‘加锁’与‘释放锁’
	// 并且无关是否需要return val，都可以用defer
}

func (c *SafeCounter) Val(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

func TestMutAdd(t *testing.T) {
	c := SafeCounter{v: make(map[string]int)}
	for range 1000 {
		go c.Inc("someKey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Val("someKey"))
}
