package basic

import (
	"fmt"
	"testing"
)

// 对于Golang来说, 实现interface的所有方法即可, 不需要明显的impl关键字
// 也就是不需要显示的声明, 这使得几乎接口的实现‘解藕’
type I interface {
	M()
	K()
}

type T struct {
	S string
}

func (t T) M() {
	fmt.Println(t.S)
}

func (t T) K() {
	fmt.Println(t.S)
}

// describe 接口可以被‘传递’, 这也是golang特殊的地方
// 接口包含: v-值, T-struct
func describe(i I) {
	i.M()
	i.K()
	fmt.Printf("%v, %T\n", i, i)
}

func TestInterfaceImpl(t *testing.T) {
	// 这里是核心, 可以var是接口I, 但是我们却使用T进行赋值,
	// 如果T没有完全实现I接口的方法, 那么这里是会编译不通过的
	// 所以对于Golang的接口来说，我们需要考虑2点:
	// 1) 如果是每个接口通用, 不需要做修改, 我们写类似describe这样接收接口的的方法即可
	// 2) 如果是每个接口调用, 各自有微调, 那就是写入接口中要实现的
	var i I = T{"Hello"}
	describe(i)
}
