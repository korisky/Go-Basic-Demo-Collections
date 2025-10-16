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

// ScaleFunc 与Scale方法做一样的内容, 但是我们并不把它'挂'在Vertex才允许调用
// 但这个时候, 就需要考虑清楚, 要将‘地址’作为接收param, 而不是普通的 Vertex2
func ScaleFunc(v *Vertex2, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// TestTheFunc 关于使用指针传递/还是使用属性方法执行的考虑
//  1. 单个struct操作内容, 那么直接作为属性方法更为合理
//  2. 如果涉及多个struct, 并且需要保留操作结果, 传递地址使用指针修改值, 会更合适
//     (但同样需要平衡性能与并发 -> 小内容直接拷贝更合理, 传递地址要避免并发同时操作)
func TestTheFunc(t *testing.T) {
	v := Vertex2{2, 3}
	fmt.Println(v.Abs())
	v.Scale(10)
	fmt.Println(v.Abs())
	// 非特定方法可调用, 我们将地址传过去, 确保func接收指针,
	// 这样也能达到类似原Scale方法的结果, 对地址其上内容进行操作
	ScaleFunc(&v, 10)
	fmt.Println(v.Abs())
	d := &v
	fmt.Println(d.Y)
}
