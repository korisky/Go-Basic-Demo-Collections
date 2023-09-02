package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type order int

func main() {

	// make a channel, user 'order' as type, which is int
	ch := make(chan order, 3)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		worker("Stringer_A", ch)
	}()

	go func() {
		defer wg.Done()
		worker("Stringer_B", ch)
	}()

	for i := 0; i < 10; i++ {
		waitForOrders()
		o := order(i)
		log.Printf("Partier: I %v, I will pass it to the channel\n", o)
		ch <- o
	}

	log.Println("No more orders, closing the channel to signify workers to stop")
	close(ch)

	log.Println("Wait for workers to gracefully stop")
	wg.Wait()

	log.Printf("All done")
}

func worker(name string, ch <-chan order) {
	// traverse the order in the channel
	for o := range ch {
		log.Printf("%s: I got %v, I will proccess it\n", name, o)
		processOrder(o)
		log.Printf("%s: Finished order %v\n", name, o)
	}
}

func waitForOrders() {
	duration := time.Duration(2+rand.Intn(2)) * time.Second
	time.Sleep(duration)
}

func processOrder(_ order) {
	duration := time.Duration(2+rand.Intn(3)) * time.Second
	time.Sleep(duration)
}
