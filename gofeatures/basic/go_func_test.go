package basic

import (
	"fmt"
	"math"
	"testing"
)

type Vertex2 struct {
	X, Y float64
}

// Abs 方法可以看到, 方法接收者v并不是指针, 这意味着对X和Y的改变
// 不会在指望完毕之后生效 (相当于执行的时候‘深拷贝’了一个Vertex2类型)
func (v Vertex2) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Scale 方法的接收者v则是指针*, 这意味着修改后, 本身v内保存的值就会改变
// 通过重复执行可以看出如此, 如果剔除*, 则会发现结果并未改变
// 带有指针*, 相当于‘将地址’拿过来然后执行
func (v *Vertex2) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func TestTheFunc(t *testing.T) {
	v := Vertex2{2, 3}
	fmt.Println(v.Abs())
	v.Scale(10)
	fmt.Println(v.Abs())
}
