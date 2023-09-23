package mistakes

import (
	"fmt"
	"testing"
)

type data struct {
	name string
}

type printer interface {
	print()
}

func (d *data) print() {
	fmt.Println("name: ", d.name)
}

func Test_pointer_negative(t *testing.T) {
	d1 := data{name: "one"}
	d1.print()

	m := map[string]data{
		"x": {"three"},
	}
	fmt.Println(m["x"].name)

	//m["x"].print() // could not call it by this under map
	theX := m["x"]
	theX.print() // but after assign it to a variable, can call the function again
}
