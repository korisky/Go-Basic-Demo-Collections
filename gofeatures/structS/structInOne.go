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

	// anonymous struct
	job := struct {
		title  string
		salary uint32
	}{
		title:  "Manager",
		salary: 100_000,
	}
	fmt.Println("Job salary: ", job.salary)
}
