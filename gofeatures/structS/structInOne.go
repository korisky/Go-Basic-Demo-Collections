package main

import (
	"encoding/json"
	"fmt"
)

type Employee struct {
	Name     string `json:"name"`
	Age      uint8  `json:"age"`
	IsRemote bool   `json:"isRemote"`
}

func main() {
	// normal struct
	employee := Employee{
		Name:     "Tom",
		Age:      32,
		IsRemote: false,
	}
	fmt.Println("Employee name: ", employee.Name)

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
	fmt.Println("New name: ", employee.Name)

	// json marshal
	jsonStr, _ := json.Marshal(employee)
	fmt.Println(string(jsonStr))
}

func (e *Employee) updateName(newName string) {
	e.Name = newName
}
