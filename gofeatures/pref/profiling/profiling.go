package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

// main: by using go-tool: go tool pprof cpu.prof
// command: web -> get the graph of the whole calling logic & cache missed/not & time spending
func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("could not create CPU profile:", err)
		return
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("could not start CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	// simulate some workload
	for i := range 10000000 {
		_ = fmt.Sprint(i)
	}
	time.Sleep(2 * time.Second)
}
