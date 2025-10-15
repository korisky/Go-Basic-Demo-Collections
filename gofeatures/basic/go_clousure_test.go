package basic

import (
	"fmt"
	"testing"
)

// adder 是一个闭包func, 作用：
// 1) sum类似全局变量, 不会被内存回收
// 2) sum变量一直存在且只有adder这个func可以触碰
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func TestClosures(t *testing.T) {
	pos, neg := adder(), adder()
	for i := range 10 {
		fmt.Println(pos(i), neg(-2*i))
	}
}
