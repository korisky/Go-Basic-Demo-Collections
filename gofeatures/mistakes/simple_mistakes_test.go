package mistakes

import (
	"fmt"
	"sort"
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

// Test_range_order, as you can see, the k & v under range would result in re-order,
// golang's compiler is intended to do this， need to take more care of the golang's behavior on 'range'
func Test_range_order(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	for k, v := range m {
		fmt.Println(k, v)
	}

	// (best way on solve this is just make your code's logic not depend on traverse order at all)
	// solution on making the order not changed at all
	// 1) sort the keys each times we update the map (or we need to traverse on it)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	fmt.Println(keys)
	sort.Strings(keys)
	fmt.Println(keys)
	// 2) use integer as our key
}
