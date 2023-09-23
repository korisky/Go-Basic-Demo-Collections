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

	//m["x"].print() // could not call it by this under map, because map can not be addressing
	theX := m["x"]
	theX.print() // but after assign it to a variable, can call the function again
}

func Test_value_changing(t *testing.T) {

	// map can not be addressing -> could not change the value directly
	m := map[string]data{
		"x": {"TOM"},
	}
	//m["x"].name = "Jerry" // could not be addressing
	d := m["x"]
	d.name = "April"
	m["x"] = d
	fmt.Println(m) // this is the way for changing the value inside map

	// but for slice, it works
	s := []data{{"Tom"}}
	s[0].name = "Jerry"
	fmt.Println(s)
}
