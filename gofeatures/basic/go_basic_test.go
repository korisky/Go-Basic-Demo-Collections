package basic

import (
	"fmt"
	"math"
	"testing"
)

// TestFuncDirectReturning 展示go允许func直接简单return的默认方式
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

// TestDefaultVals 展示不同类型默认值
func TestDefaultVals(t *testing.T) {
	var i int // 不需要+-的时候多使用uint, 与Java不同, Go支持默认无符号会使得与Solana交互非常方便
	var f float64
	var b bool
	var s string
	// 这里使用%s的情况, 空字符串直接是空, 但用%q会展示""
	fmt.Printf("%v, %v, %v, %q\n", i, f, b, s)
}

// TestHardTypeConvert 与Java不同,
// Go的一切类型转换都需要显示指明, 但不需要区分封装类型或者其他的, 直接使用类型方法进行转换即可
func TestHardTypeConvert(t *testing.T) {
	var x, y int = 3, 4
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z)
}
