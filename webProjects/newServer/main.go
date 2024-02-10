package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
func main() {

	// main function only calls run()
	ctx := context.Background()
	// can be cancelled by Ctrl + C
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}

}

// run function's last input param -> args -> is for environment variables
func run(ctx context.Context, w io.Writer, args []string) error {
	// ...
	return nil
}

// handleSomething letting the function return handler itself,
func handleSomething(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// logging
		})
}
