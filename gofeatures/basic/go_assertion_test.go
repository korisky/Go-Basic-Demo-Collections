package basic

import (
	"fmt"
	"testing"
)

// TestAssertion Golang中的断言(assertion)可以提供转换判断
func TestAssertion(t *testing.T) {

	// 初始化hello字符串, 通过断言, 哪怕使用的最泛的空接口 interface{}
	// 我们也有判断类型的方式, ok直接知道是否合理
	var i interface{} = "hello"
	s, ok := i.(string)
	fmt.Println(s, ok)

	// 这里将会是false因为不属于float64类型
	f, ok := i.(float64)
	fmt.Println(f, ok)

	// 断言失败, 还在强行赋值, 那么会直接给一个panic
	f = i.(float64)
	fmt.Println(f)
}
