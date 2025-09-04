package goroutines_test

import (
	"fmt"
	"testing"
	"time"
)

func TestSimpleRoutine(t *testing.T) {
	go say("world")
	say("hello")
}

func say(s string) {
	for range 5 {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
