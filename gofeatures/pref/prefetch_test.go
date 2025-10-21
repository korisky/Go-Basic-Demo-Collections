package pref

type Node struct {
	Value int
	Data  [64]byte
}

func SumLinear(data []int) int {
	sum := 0
	for i := 0; i < len(data); i++ {
		// 由于我们是连续获取, 所以CPU可以prefetch下一个到后面的
		sum += data[i]
	}
	return sum
}

func SunRandom(data []int, indices []int) int {
	sum := 0
	for _, idx := range indices {
		// 由于我们是随意获取, cache miss
		sum += data[idx]
	}
	return sum
}

// ProcessWithPrefetch 我们使用手动prefetch
//func ProcessWithPrefetch(nodes []*Node) int {
//	sum := 0
//	for i := 0; i < len(nodes)-4; i++ {
//		// 这里就是特殊操作, _说明我不使用, 纯加载
//		_ = nodes[i + 4].Value
//
//		sum +=
//	}
//}
