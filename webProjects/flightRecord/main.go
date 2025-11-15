package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/trace"
	"sync"
	"time"
)

func heavyLoad(w *sync.WaitGroup, iterations int32) {
	defer w.Done()
	for i := range iterations {
		_ = fmt.Sprintf("processing %d", i)
	}
	time.Sleep(500 * time.Millisecond)
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		var wg sync.WaitGroup
		wg.Add(2)

		go heavyLoad(&wg, 100_000)
		go heavyLoad(&wg, 10_000_000)

		wg.Wait()

		consume := time.Since(start)
		_, err := fmt.Fprintf(w, "worked for %f seconds", consume.Seconds())
		if err != nil {
			return
		}
	}
}

// main is for simulation for checking the FlightRecord
// Flight Record 使用类似时间轮的方式, 记录最近几次记录的tracing
// 所以定义时间窗口的大小, 是使用的关键.
func main() {

	// flight-record config, 记住MAXBytes永远为上限, 一旦达到则直接记录即使没有达到配置的MinAge¬
	fr_cfg := trace.FlightRecorderConfig{
		MinAge:   5 * time.Second, // 定义fc的窗口lower-bound
		MaxBytes: 3 << 20,         // 3MB, 1 << 10: 1kb, 1 << 20 Mb, 1 << 30: 1Gb
	}

	// start flight-record
	fr := trace.NewFlightRecorder(fr_cfg)
	err := fr.Start()
	if err != nil {
		log.Fatalf("unalbe to start trace flight recorder: %v", err)
	}
	defer fr.Stop()

	// start serving
	http.HandleFunc("/", handler())
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
