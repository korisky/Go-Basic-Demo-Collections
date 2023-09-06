package main

import (
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func main() {

	keyFile := "/Users/roylic/PEM/key.unencrypted.pem"
	certFile := "/Users/roylic/PEM/cacert.pem"

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, Http3"))
	})
	err := http3.ListenAndServeQUIC("localhost:4242", certFile, keyFile, nil)
	if err != nil {
		panic(err)
	}
}
