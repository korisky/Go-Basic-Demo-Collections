package main

import (
	"fmt"
	"log"
	"net/http"
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
func main() {
	http.HandleFunc("/", handler())
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
