package streamActions

import "testing"

func Test_ChainFrom(t *testing.T) {
	From([]int{1, 2, 3, 4}).Each(func(a int) { println(a) })
}

func Test_Reverse(t *testing.T) {
	From([]int{1, 2, 3, 4}).Reverse().Each(func(a int) { println(a) })
}
