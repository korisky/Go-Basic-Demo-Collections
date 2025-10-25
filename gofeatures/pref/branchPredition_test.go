package pref

import (
	"math/rand"
	"sort"
	"testing"
)

const theDataSize = 10000000

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

func setupRandomData() []int {
	data := make([]int, theDataSize)
	r := rand.New(rand.NewSource(42))
	for i := range data {
		data[i] = r.Intn(256)
	}
	return data
}

func setupSortedData() []int {
	data := setupRandomData()
	sort.Ints(data)
	return data
}

func BenchmarkForBranchPrediction(b *testing.B) {
	b.Run("Unpredictable-Random", func(b *testing.B) {
		data := setupRandomData()
		b.ResetTimer()
		for range b.N {
			_ = CountConditionRandom(data, 128)
		}
	})
	b.Run("Predictable-PreSort", func(b *testing.B) {
		data := setupSortedData()
		b.ResetTimer()
		for range b.N {
			_ = CountConditionRandom(data, 128)
		}
	})
	b.Run("Sorted-WithSortCost", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			data := setupRandomData()
			_ = CountConditionSorted(data, 128)
		}
	})
	b.Run("Branchless-General", func(b *testing.B) {
		data := setupRandomData()
		b.ResetTimer()
		for range b.N {
			_ = CountConditionBranchless(data, 128)
		}
	})
}
