package basic

import (
	"fmt"
	"testing"
)

// Index 用于进行查找
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// 找到则返回
		if v == x {
			return i
		}
	}
	return -1
}

func TestComparable(t *testing.T) {
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))
	fmt.Println(Index(si, -105))
	fmt.Println()

	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))
	fmt.Println(Index(ss, "baz"))
}
