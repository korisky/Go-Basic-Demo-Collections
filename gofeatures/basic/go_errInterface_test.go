package basic

import (
	"fmt"
	"testing"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v,\n got error with reason:%s", e.When, e.What)
}

func (e *MyError) RuntimeError() {
	panic("Runtime Error Occur")
}

func run() error {
	return &MyError{
		When: time.Now(),
		What: "Just an error",
	}
}

func TestErrImpl(t *testing.T) {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
