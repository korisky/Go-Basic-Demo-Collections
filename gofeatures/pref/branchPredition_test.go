package pref

import "sort"

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
		// 2. >> 63 获取最高位, 1正数, 0负数
		// 3. 与1进行异或XOR操作, 1也就是正数也就添加,
		count += int(uint(v-threshold-1)>>63) ^ 1
	}
	return count
}
