package basic

import (
	"fmt"
	"testing"
)

func TestFuncDirectReturning(t *testing.T) {
	fmt.Println(split(35))
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	// 对于go语言来说, 如果恰好使用同名的返回参数
	// 直接进行return会默认返回操作值
	// 等同于 return x, y
	return
}
