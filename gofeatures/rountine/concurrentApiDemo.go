package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/products",
		"https://api.example.com/orders",
	}

	results := make(chan string)

	for _, url := range urls {
		go fetchAPI(ctx)
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
