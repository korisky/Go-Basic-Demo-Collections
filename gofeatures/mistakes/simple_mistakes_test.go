package mistakes

import (
	"fmt"
	"testing"
)

// Test_closure, when we use := rather than =, it only works in {}, but = could work for everywhere
func Test_closure(t *testing.T) {
	x := 1
	y := 10
	fmt.Printf("X:%v, Y:%v\n", x, y)
	{
		fmt.Printf("X:%v, Y:%v\n", x, y)
		x := 2
		y = 20
		fmt.Printf("X:%v, Y:%v\n", x, y)
	}
	fmt.Printf("X:%v, Y:%v\n", x, y)
}

func Test_closure_v2(t *testing.T) {
	x := [3]int{1, 2, 3}
	func(arr [3]int) {
		// In golang, array also wold send 'value copy', unlike Java & C++, they send array as 'reference'
		// thus, in golang, array's value manipulation would not remain after stack exit
		arr[0] = 7
		fmt.Println(arr)
	}(x)
	fmt.Println(x)
	fmt.Println()

	// two saving ways ->
	// 1) pass pointer, and manipulate the value
	y := [3]int{1, 2, 3}
	func(arr *[3]int) {
		arr[0] = 7
		fmt.Println(*arr)
	}(&y)
	fmt.Println(y)
	fmt.Println()

	// 2) use slice rather than array
	z := y[:]
	func(arr []int) {
		arr[0] = 9
		fmt.Println(arr)
	}(z)
	fmt.Println(y)
	fmt.Println(z)
	fmt.Println()
}

func Test_map_contains(t *testing.T) {

	x := map[string]string{"one": "1", "two": ""}

	// error on checking a key is exist or not in a map
	if v := x["two"]; v == "" {
		fmt.Println("key two is not exist")
	}

	// we should check it by using 2 output param
	if _, exist := x["two"]; !exist {
		fmt.Println("key two is not exist")
	}

}
