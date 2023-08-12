package gc

import (
	"fmt"
	"runtime"
	"testing"
)

// In golang, using a slice as buffer -> better then
func Test_GCStuff(t *testing.T) {

	s := make([]string, 0, 100000)
	for i := 0; i < 100000; i++ {
		s = append(s, "hello, world")
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("HeapAlloc: ", m.HeapAlloc)
	fmt.Println("HeapIdle: ", m.HeapIdle)
	fmt.Println("HeapReleased: ", m.HeapReleased)
	fmt.Println("NumGC: ", m.NumGC)
	fmt.Println("-----------")

	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println("HeapAlloc: ", m.HeapAlloc)
	fmt.Println("HeapIdle: ", m.HeapIdle)
	fmt.Println("HeapReleased: ", m.HeapReleased)
	fmt.Println("NumGC: ", m.NumGC)
	fmt.Println("-----------")

	// when a variable is no longer need to be used,
	// set it to nil allows the gc to reclaim the memory from it
	s = nil
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println("HeapAlloc: ", m.HeapAlloc)
	fmt.Println("HeapIdle: ", m.HeapIdle)
	fmt.Println("HeapReleased: ", m.HeapReleased)
	fmt.Println("NumGC: ", m.NumGC)
	fmt.Println("-----------")
}
