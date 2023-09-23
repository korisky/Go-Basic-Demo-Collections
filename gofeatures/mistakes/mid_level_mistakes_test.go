package mistakes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_http_close(t *testing.T) {
	// first check the request error
	resp, err := http.Get("https://www.google.com")
	if resp != nil {
		// is resp is not nil, do not forget the close it later
		// if forget to close it (e.g. when it return nil), the whole program would panic and stop
		defer resp.Body.Close()
	}
	checkError(err)

	body, err := io.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func Test_recover(t *testing.T) {
	defer func() {
		fmt.Println("recovered: ", recover())
	}()
	panic("Not good")
}

func Test_loop_and_goroutine(t *testing.T) {

	/*
		Detailed explanation:
			- the goroutines would be executed after scheduling, not exactly the time when code find 'go func....'
			- each time meet the 'go func...', they would record params they need. In below example, all goroutines catch
			the variable v. But only it's reference (record the pointer for this variable v)
			- so, it related to keyword 'range', it would treat v as everyone inside the loop only get 'pointer' of v
			- thus, when scheduler let the goroutine to be executed, the loop already finished, v is updated from one -> two -> three
			and go func() only get 'three'
			- for avoiding this, we can simply pass the variable (copy as calling), or init a new local variable, could solve this
	*/

	// no matter you are traversing over slice, array, or even map
	data := []string{"one", "two", "three"}

	for _, v := range data {
		// because of closure in Golang, all goroutines would result in getting same variable v
		// when the iteration is stop, thus, all print THREE rather than ONE, TWO, THREE

		vCopy := v // by adding vCopy (or we can pass the value into the anonymous function), we could solve the problem
		go func() {
			fmt.Println(vCopy)
		}()
	}
	time.Sleep(10 * time.Second)
	// out-print three three three
}

// Test_defer_param for the defer function, the param of it, would be settled on it's declaration
func Test_defer_param(t *testing.T) {
	var i = 1
	// for defer function, if we use it inside the for loop, it would be closed,
	// it will only be executed when the main thread the called is going to be terminated
	defer fmt.Println("result: ", func() int { return i * 5 }())
	i++ // golang only support i++, could not ++i
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
