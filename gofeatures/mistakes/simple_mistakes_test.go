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
