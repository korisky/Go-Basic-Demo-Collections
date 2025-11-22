package main

import "fmt"

type Employee struct {
	name     string
	age      uint8
	isRemote bool
}

func main() {
	// normal struct
	employee := Employee{
		name:     "Winston",
		age:      32,
		isRemote: false,
	}
	fmt.Println("Employee name: ", employee.name)
}
