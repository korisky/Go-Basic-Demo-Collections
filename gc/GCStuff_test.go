package gc

import (
	"fmt"
	"runtime"
	"testing"
)

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

	s = nil
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println("HeapAlloc: ", m.HeapAlloc)
	fmt.Println("HeapIdle: ", m.HeapIdle)
	fmt.Println("HeapReleased: ", m.HeapReleased)
	fmt.Println("NumGC: ", m.NumGC)
	fmt.Println("-----------")
}
