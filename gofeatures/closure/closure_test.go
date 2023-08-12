package closure

import (
	"fmt"
	"testing"
)

// Test_Closure is an example for closure feature in golang
// compare to Java's anonymous inner class, which can only access final variables
// golang's closure can be used with much more flexibility
func Test_Closure(t *testing.T) {
	x := 10

	// -> it might need to promotes local variable x from stack to heap -> when it met it in closure
	increment := func() int {
		x++
		return x
	}
	fmt.Println(increment())
	fmt.Println(increment())
}

// Test_ClosureCallback usually we use closure in golang as we manipulate Callback
func Test_ClosureCallback(t *testing.T) {
	f := adder(10)
	fmt.Println(f(5))
}

func adder(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}
