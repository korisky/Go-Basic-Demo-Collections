package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	// create a context -> with value can let the context 'hold' values
	// similar to Java's thread_local?
	ctx := context.Background()
	ctx = context.WithValue(ctx, "UserId", 10086)

	go performPassing(ctx, &wg)

	wg.Wait()
}

// performPassing
func performPassing(ctx context.Context, wg *sync.WaitGroup) {
	// get value from the passing context
	userID := ctx.Value("UserId")
	fmt.Println("User ID:", userID)
	wg.Done()
}
