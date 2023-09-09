package main

import (
	"context"
	"fmt"
	"time"
)

// below example -> construct a context with 2s timeout -> try to do a 5-second-duration task
// then the timeout ctx would automatically close (.Done()) when it reach 2s
func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	go performTask(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("Task exceed ctx timeout")
	}
}

// performTask simulate a time-consuming task(5 seconds)
func performTask(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Task completed successfully")
	}
}
