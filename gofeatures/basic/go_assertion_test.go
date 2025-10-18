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

// selectByType 空接口interface{} + select 判断类型
// 可以做到用最泛的接口（因为不需要实现）, 然后重新回退到具体类型进行处理
func selectByType(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer Detected: %v\n", v)
	case string:
		fmt.Printf("String Detected: %v\n", v)
	default:
		fmt.Printf("What type is it? %T\n", v)
	}
}

func TestAssertionWithSelect(t *testing.T) {
	selectByType(30)
	selectByType("abc")
	selectByType(true)
}
