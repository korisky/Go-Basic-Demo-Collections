package goroutines_test

import (
	"fmt"
	"testing"
)

func TestChanStuff(t *testing.T) {

	// init the array
	s := []int{7, 2, 8, -9, 4, 0, 5, 12, 3354, 1234, 52, -2523, 123}
	c := make(chan int)

	// async summing
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	// retrieve from chan
	x, y := <-c, <-c

	fmt.Println(x, y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send result to channel
}
