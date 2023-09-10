package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	// context that could be canceled
	ctx, cancelFunc := context.WithCancel(context.Background())
	go performCancel(ctx)
	time.Sleep(2 * time.Second)

	// call the cancel func -> trigger to context, the context would received -> .Done()
	cancelFunc()
	time.Sleep(5 * time.Second)
}

func performCancel(ctx context.Context) {
	for {
		select {
		// if ctx is done -> must be called cancelled from outside
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			fmt.Println("Performing task...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
