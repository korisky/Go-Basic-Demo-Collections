package basic

import (
	"fmt"
	"testing"
)

// Index 用于进行查找, 与Java不同, golang中的默认类型是支持直接进行==比较, 这里比较的就是值
// 但对于Golang中的ref类型 (map, slice, func, struct), 那还是需要实现Comparable接口
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
