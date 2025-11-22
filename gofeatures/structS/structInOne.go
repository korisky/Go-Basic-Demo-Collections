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
		name:     "Tom",
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

	// pointer -> ref of struct
	employee.updateName("Winston")
	fmt.Println("New name: ", employee.name)
}

func (e *Employee) updateName(newName string) {
	e.name = newName
}
