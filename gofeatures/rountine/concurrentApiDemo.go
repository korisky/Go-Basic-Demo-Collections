package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// The http.NewRequestWithContext() function is used to create an HTTP request with the provided context.
// If any of the API requests exceed the timeout duration, the context's cancellation signal is propagated,
// canceling all other ongoing requests.
func main() {
	// create a context with timeout func
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/products",
		"https://api.example.com/orders",
	}
	// create a channel
	results := make(chan string)
	// concurrently do the request -> new go routine
	for _, url := range urls {
		// same ctx means: one of them is timeout -> all of them till be canceled
		go fetchAPI(ctx, url, results)
	}
	// print results
	for range urls {
		fmt.Println(<-results)
	}
}

func fetchAPI(ctx context.Context, url string, results chan<- string) {
	// build request with context (context might contains timeout stuff)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		results <- fmt.Sprintf("Error creating request for %s: %s", url, err.Error())
		return
	}
	// make the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		results <- fmt.Sprintf("Error making request for %s: %s", url, err.Error())
		return
	}
	// read the response
	defer resp.Body.Close()
	results <- fmt.Sprintf("Response from %s: %d", url, resp.StatusCode)
}
