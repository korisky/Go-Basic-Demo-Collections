package closure

import (
	"fmt"
	"testing"
)

// Test_Closure_error showing that, in Golang, the anonymous function inside the
// for loop, captures the variable 'v' by REFERENCE, thus, since the for loop
// update 'v' with each iteration, all the anonymous functions end up sharing the same
// reference to 'v'
func Test_Closure_error(t *testing.T) {

	done := make(chan bool)
	values := []string{"a", "b", "c"}

	for _, v := range values {
		go func() {
			fmt.Println(v)
			done <- true // trigger finished
		}()
	}

	for _ = range values {
		<-done // wait execution
	}
}
