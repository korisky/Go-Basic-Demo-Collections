package test_map

import (
	"fmt"
	"math"
	"testing"
)

func Test_SumArray(t *testing.T) {
	arr := [10]int{2, 5, 6, 7, 87, 98}
	fmt.Println(SumArray(arr[:])) // 由于传输的slice, 使用[:]把array转一个slice传入
	fmt.Println(MaxArray(arr[:])) // 由于传输的slice, 使用[:]把array转一个slice传入

	fmt.Printf("%+v\n", WordCount("bananananana"))
	fmt.Printf("%+v\n", FilterTrans(arr[:]))

	s1, s2, s3, s4 := "abc", "bcs", "225s", "1235sdf"
	sArr := []string{s1, s2, s3, s4}

	fmt.Printf("%+v\n", GroupByLen(sArr[:]))
}

func SumArray(arr []int) int {
	sum := 0
	for _, val := range arr {
		sum += val
	}
	return sum
}

func MaxArray(arr []int) int {
	theMax := math.MinInt
	for _, val := range arr {
		theMax = max(theMax, val)
	}
	return theMax
}

func WordCount(str string) map[string]int {
	freqMap := make(map[string]int)
	for _, val := range str {
		freqMap[string(val)]++
	}
	return freqMap
}

func FilterTrans(arr []int) []int {
	var target []int
	for _, val := range arr {
		if val%2 == 1 {
			target = append(target, val*2)
		}
	}
	return target
}

func GroupByLen(strArr []string) map[int][]string {
	res := make(map[int][]string)
	for _, str := range strArr {
		strLen := len(str)
		res[strLen] = append(res[strLen], str)
	}
	return res
}
