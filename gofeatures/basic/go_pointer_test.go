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

type Vertex struct {
	X int
	Y int
}

func TestStructWithPointer(t *testing.T) {
	v := Vertex{1, 2}
	p := &v // p是指向v的指针
	p.X = 1e9
	fmt.Println(v)
	// 可以发现, 当作指针*获取指向内容的属性X, 与直接p.X是一样的
	fmt.Println((*p).X)
	fmt.Println(p.X)
}

// TestArraySlice 区分数组(固定)与切片(动态)
func TestArraySlice(t *testing.T) {
	arr := [6]int{1, 2, 3, 4, 5}
	s1 := arr[0:3]
	s2 := arr[2:]
	fmt.Println(arr)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println()

	// arr变化对于切片的影响
	arr[2] = 100
	fmt.Println(arr)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println()

	// 不管是修改原来的arr, 还是只修改slice
	// 都意味着在原arr进行了修改，所有相关的都生效
	s1[2] = 98
	fmt.Println(arr)
	fmt.Println(s1)
	fmt.Println(s2)
}
