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

func TestInterfaceImpl(t *testing.T) {
	var i I = T{"Hello"}
	i.M()
	i.K()
}
