package pref

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

// CacheLineSize 可以作为环境变量, 直接放入文件使用
const CacheLineSize = 128

// Counters interface -> better reuse later benchmark code
type Counters interface {
	RequestsPtr() *uint64
	ErrorPtr() *uint64
	LatencyPtr() *uint64
}

// CountersWithoutPadding will cause 'false-sharing'
type CountersWithoutPadding struct {
	requests uint64
	errors   uint64
	latency  uint64
}

func (c *CountersWithoutPadding) RequestsPtr() *uint64 {
	return &c.requests
}

func (c *CountersWithoutPadding) ErrorPtr() *uint64 {
	return &c.errors
}

func (c *CountersWithoutPadding) LatencyPtr() *uint64 {
	return &c.latency
}

// CountersWithPadding since with extra padding, fill the whole cache line
// then it will not cause 'false-sharing', and result in a greater performance
type CountersWithPadding struct {
	requests uint64
	_        [CacheLineSize - 8]byte
	errors   uint64
	_        [CacheLineSize - 8]byte
	latency  uint64
	_        [CacheLineSize - 8]byte
}

func (c *CountersWithPadding) RequestsPtr() *uint64 {
	return &c.requests
}

func (c *CountersWithPadding) ErrorPtr() *uint64 {
	return &c.errors
}

func (c *CountersWithPadding) LatencyPtr() *uint64 {
	return &c.latency
}

// runBenchMark create multiple workers to do benchmark
// improvement -> 使用 b.Loop 而不是 b.N , 强制将构建goroutine的耗费部分在循环外，
// 更准确的benchmark出我们写的func性能, 剔除了goroutine生成的耗时
func runBenchMark(b *testing.B, counters Counters) {

	var wg sync.WaitGroup
	stop := make(chan struct{})

	b.ResetTimer()

	// Spawn fixed workers ONCE
	// 一次声明所有的worker要做的func
	workers := []func(){
		func() {
			ptr := counters.RequestsPtr()
			for {
				select {
				case <-stop:
					return
				default:
					atomic.AddUint64(ptr, 1)
				}
			}
		},
		func() {
			etr := counters.ErrorPtr()
			for {
				select {
				case <-stop:
					return
				default:
					atomic.AddUint64(etr, 1)
				}
			}
		},
		func() {
			ltr := counters.LatencyPtr()
			for {
				select {
				case <-stop:
					return
				default:
					atomic.AddUint64(ltr, 1)
				}
			}
		},
	}

	// 在这个地方统一分配go-routine
	for _, w := range workers {
		wg.Add(1)
		// worker由于带有自己的func, 针对传入func作为处理内容
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(w)
	}

	b.ResetTimer()

	// 通过使用Loop, 也可以避免很多compiler-optimized内容, 更直接的验证func
	for b.Loop() {

	}

	close(stop)
	wg.Wait()
}

// BenchmarkComparison benchmarking one-one comparison
func BenchmarkComparison(b *testing.B) {
	b.Run("NoPadding-FalseSharing", func(b *testing.B) {
		var counters CountersWithoutPadding
		runBenchMark(b, &counters)
	})
	fmt.Println()
	b.Run("WithPadding-NoFalseSharing", func(b *testing.B) {
		var counters CountersWithPadding
		runBenchMark(b, &counters)
	})
}

func TestMemLayout(t *testing.T) {
	// size using show
	unPaddedSize := unsafe.Sizeof(CountersWithoutPadding{})
	paddedSize := unsafe.Sizeof(CountersWithPadding{})
	t.Logf("CounterWithoutPadding size: %d bytes (all in 1 cache-line)", unPaddedSize)
	t.Logf("CounterWithPadding size: %d bytes (3 separate cache-line)", paddedSize)

	// size using check
	if unPaddedSize >= 128 {
		t.Errorf("CounterWithoutPadding too large: %d bytes (should be < 128, the CPU single core cache-line size)", paddedSize)
	}
	expectedPaddedSize := uint(128 * 3)
	if paddedSize != uintptr(expectedPaddedSize) {
		t.Errorf("CounterWithPadding size mis-match: got %d bytes, want %d", paddedSize, expectedPaddedSize)
	}
}
