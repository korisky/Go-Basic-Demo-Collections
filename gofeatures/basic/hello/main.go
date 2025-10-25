package main

import (
	"fmt"
	"own/simple/greetings"
)

func main() {
	msg := greetings.Hello("Gladys")
	fmt.Print(msg)
}
