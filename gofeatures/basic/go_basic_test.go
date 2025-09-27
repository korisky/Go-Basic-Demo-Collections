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

// 常量的生命可以在后面+具体类型, 达到控制精度等效果
const TheConst uint32 = 15

// TestDefaultVals 展示不同类型默认值
func TestDefaultVals(t *testing.T) {
	var i int // 不需要+-的时候多使用uint, 与Java不同, Go支持默认无符号会使得与Solana交互非常方便
	var f float64
	var b bool
	var s string
	// 这里使用%s的情况, 空字符串直接是空, 但用%q会展示""
	// %T直接展示类型
	fmt.Printf("%v is type %T\n%v is type %T\n%v is type %T\n%q is type %T\n%v is type %T\n",
		i, i, f, f, b, b, s, s, TheConst, TheConst)
}

// TestHardTypeConvert 与Java不同,
// Go的一切类型转换都需要显示指明, 但不需要区分封装类型或者其他的, 直接使用类型方法进行转换即可
func TestHardTypeConvert(t *testing.T) {
	var x, y int = 3, 4
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z)
}

// TestIfWithVar 展示了Go中对判断结果进一步判断走入branch的简洁写法
func TestIfWithVar(t *testing.T) {
	fmt.Println(
		pow(3, 2, 10),
		pow(3, 3, 20),
	)
}

func pow(x, n, lim float64) float64 {
	// 判断得到的值, 进行二次判断, 从而使得多次判断代码更简洁
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}
