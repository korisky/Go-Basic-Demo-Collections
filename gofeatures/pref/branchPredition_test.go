package pref

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func CountConditionRandom(data []int, threshold int) int {
	count := 0
	for _, v := range data {
		if v > threshold {
			count++
		}
	}
	return count
}

func CountConditionSorted(data []int, threshold int) int {
	sort.Ints(data)
	count := 0
	for _, v := range data {
		if v > threshold {
			count++
		}
	}
	return count
}

func CountConditionBranchless(data []int, threshold int) int {
	count := 0
	for _, v := range data {
		// 1. 如果 v - threshold - 1 < 0, 我们主要看结果是正数还是负数
		// 2. >> 63 获取最高位, 1负数, 0正数
		// 3. 与1进行异或XOR操作, 1也就是负数也就不添加, 0是正数则添加
		count += int(uint(v-threshold-1)>>63) ^ 1
	}
	return count
}

// CountConditionCMOV 使用到汇编中的CMOV指令进行处理
func CountConditionCMOV(data []int, threshold int) int {
	count := 0
	for _, v := range data {
		increment := 0
		if v > threshold {
			increment = 1
		}
		count += increment
	}
	return count
}

func setupRandomData() []int {
	data := make([]int, theDataSize)
	r := rand.New(rand.NewSource(42))
	for i := range data {
		data[i] = r.Intn(1024)
	}
	return data
}

func setupSortedData() []int {
	data := setupRandomData()
	sort.Ints(data)
	return data
}

const theDataSize = 100000000
const threshold = 512

func BenchmarkForBranchPrediction(b *testing.B) {

	randomData := setupRandomData()
	sortedData := setupSortedData()

	b.Run("Unpredictable-Random", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_ = CountConditionRandom(randomData, threshold)
		}
	})
	fmt.Println()

	b.Run("Predictable-PreSort", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_ = CountConditionRandom(sortedData, threshold)
		}
	})
	fmt.Println()

	b.Run("Sorted-WithSortCost", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			// with sort cost inner
			data := append([]int(nil), randomData...)
			sort.Ints(data)
			_ = CountConditionSorted(data, threshold)
		}
	})
	fmt.Println()

	b.Run("Branchless-General", func(b *testing.B) {
		data := setupRandomData()
		b.ResetTimer()
		for range b.N {
			_ = CountConditionBranchless(data, threshold)
		}
	})
	fmt.Println()

	b.Run("05-Branchless-Sorted", func(b *testing.B) {
		// Show branchless is consistent regardless of data pattern
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CountConditionBranchless(sortedData, threshold)
		}
	})
	fmt.Println()

	b.Run("06-CMOV-Random", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CountConditionCMOV(randomData, threshold)
		}
	})
}
