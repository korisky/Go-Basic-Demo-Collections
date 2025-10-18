package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

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
	_        [120]byte
	errors   uint64
	_        [120]byte
	latency  uint64
	_        [120]byte
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
func runBenchMark(b *testing.B, counters Counters) {

	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			ptr := counters.RequestsPtr()
			for range 1000 {
				atomic.AddUint64(ptr, 1)
			}
		}()
		go func() {
			defer wg.Done()
			etr := counters.ErrorPtr()
			for range 1000 {
				atomic.AddUint64(etr, 1)
			}
		}()
		go func() {
			defer wg.Done()
			ltr := counters.LatencyPtr()
			for range 1000 {
				atomic.AddUint64(ltr, 1)
			}
		}()

		wg.Wait()
	}
}

// BenchmarkHasFalseSharing 不使用padding, 可以看出最大也就是到 24176
func BenchmarkHasFalseSharing(b *testing.B) {
	var counters CountersWithoutPadding
	runBenchMark(b, &counters)
}

// BenchmarkHasNoFalseSharing 使用padding, 能到 55000, 接近not-padding的2倍
func BenchmarkHasNoFalseSharing(b *testing.B) {
	var counters CountersWithPadding
	runBenchMark(b, &counters)
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
