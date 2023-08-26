package storage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"own/example/bolbol/entity"
	"runtime"
	"testing"
)

func BenchmarkMemoryWithChannel_PushNewItem(b *testing.B) {
	benchmarkMemory_PushNewItem(NewMemoryWithChannel(1000), b)
}

func BenchmarkMemoryWithList_PushNewItem(b *testing.B) {
	benchmarkMemory_PushNewItem(NewMemoryWithList(1000), b)
}

func testNewMemory(m Storage, t *testing.T) {
	ctx := context.Background()
	_ = m.Push(ctx, 10, entity.UnreadMessagesNotification{Count: 1})
	_ = m.Push(ctx, 10, entity.UnreadMessagesNotification{Count: 2})
	_ = m.Push(ctx, 10, entity.UnreadMessagesNotification{Count: 3})
	c, _ := m.Count(ctx, 10)
	assert.Equal(t, 3, c)

	p, err := m.Pop(ctx, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, p.(entity.UnreadMessagesNotification).Count)

	all, _ := m.PopAll(ctx, 10)
	assert.Equal(t, 2, len(all))

	for i := 0; i < 15; i++ {
		_ = m.Push(ctx, 10, entity.UnreadMessagesNotification{Count: i})
	}

	f, err := m.Pop(ctx, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, f.(entity.UnreadMessagesNotification).Count)
}

func benchmarkMemory_PushAverage(m Storage, b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		id := rand.Intn(1000)
		_ = m.Push(ctx, id, entity.UnreadMessagesNotification{Count: i})
	}
	b.StopTimer()
	PrintMemUsage()
}

func benchmarkMemory_PushNewItem(m Storage, b *testing.B) {
	ctx := context.Background()
	counter := 0
	for i := 0; i < b.N; i++ {
		_ = m.Push(ctx, i, entity.UnreadMessagesNotification{Count: i})
		counter++
	}
	b.StopTimer()
	b.Log("for ", b.N, " notifications: ")
	PrintMemUsage()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
