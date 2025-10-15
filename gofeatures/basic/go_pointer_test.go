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

func TestSliceBasic(t *testing.T) {
	s := []int{2, 3, 4, 5, 6, 7}
	printSlice(s)

	// 清除所有内容
	s = s[:0]
	printSlice(s) // len=0，cap仍为6

	// 拓展长度
	s = s[:4]
	printSlice(s) // len=4, cap仍为6

	// 剔除前2个值, 也就是让slice从2开始
	s = s[2:]
	printSlice(s) // len=2，cap减少为4

	// 这个时候如果我们回去, 也无法回去, 开始位置仍然从4开始, 所以cap也不会变大
	// 所以在进行slice的裁剪操作时, 尤其是没有记录array的引用, 是有可能丢失数据
	// 虽然仍存在mem中，但go并没有提供方法直接访问
	s = s[0:]
	printSlice(s) // Still len=2, cap=4, [4, 5]
}

func printSlice(s []int) {
	fmt.Printf("len=%d, cap=%d, %v\n", len(s), cap(s), s)
}

// TestSliceMake -> 容量始终是从切片的当前起始指针到底层数组(声明的/默认的)的末尾计算的
// 1. 没有指定cap的情况， 默认cap=len
// 2. 当没有arr的情况构建slice (e.g make([]int), 这样是声明slice所以没有arr)
// 3. 当进行cut-slice的时候, 当前len就是cut出来的长度, 但cap取决于:
// 3.1) 如果cut是从头开始, 那len是=cut出来长度, 但cap仍为原来大小
// 3.2) 如果cut是中间开始, 那么len是=cut出来的长度, 但cap则是从开始位置到结束
func TestSliceMake(t *testing.T) {
	a := make([]int, 5)    // len(5), then cap default is 5
	b := make([]int, 0, 5) // len(0), but cap need 5
	c := b[:2]             // len(2), but cap still 5 !!!
	d := c[2:5]            // len(3), but cap now 3 !!!
	e := d[1:]             // len(2), cap(2)
	printSlice(a)
	printSlice(b)
	printSlice(c)
	printSlice(d)
	printSlice(e)
}
