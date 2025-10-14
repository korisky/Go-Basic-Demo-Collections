package basic

import (
	"fmt"
	"testing"
)

func TestPointer(t *testing.T) {

	// init
	i, j := 42, 2701

	// p是指向i的指针
	p := &i
	fmt.Println(*p) // &获取的是i的值，*获取指向的内容
	fmt.Println(p)  // p可以看出是0x开头16进制的地址信息

	// 通过指针, 修改i的值
	*p = 21
	fmt.Println(*p) // *获取指向的内容, 发现已经出现了变化
	fmt.Println(i)  // 再看i值, 确实发生了变化

	// p转而指向j
	p = &j
	*p = *p / 37
	fmt.Println(i) // 查看i的变化
	fmt.Println(j) // 查看j的变化
}
