package main

import (
	"fmt"
	"log"
	"own/simple/greetings"
)

func main() {
	// predefined logger
	log.SetPrefix("log-greetings: ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	msg, err := greetings.Hello("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(msg)
}
