package main

import (
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func main() {
	client := http.Client{Transport: &http3.RoundTripper{}}
}
