package pref

import (
	"math/rand"
	"testing"
)

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

func SumRandom(data []int, indices []int) int {
	sum := 0
	for _, idx := range indices {
		// 由于我们是随意获取, cache miss
		sum += data[idx]
	}
	return sum
}

// ProcessWithPrefetch 我们使用手动prefetch
func ProcessWithPrefetch(nodes []*Node) int {
	sum := 0
	for i := 0; i < len(nodes)-4; i++ {
		// 这里就是特殊操作, _说明我不使用, 纯加载
		_ = nodes[i+4].Value
		// prefetch
		sum += expensiveOperation(nodes[i])
	}
	// last 4
	for i := len(nodes) - 4; i < len(nodes); i++ {
		sum += expensiveOperation(nodes[i])
	}
	return sum
}

func ProcessWithoutPrefetch(nodes []*Node) int {
	sum := 0
	for i := 0; i < len(nodes); i++ {
		sum += expensiveOperation(nodes[i])
	}
	return sum
}

// expensiveOperation 模拟高负载操作
func expensiveOperation(node *Node) int {
	result := node.Value
	for i := 0; i < 10; i++ {
		result = result*31 + int(node.Data[i])
	}
	return result
}

func setupLinearData() []int {
	data := make([]int, dataSize)
	for i := range data {
		data[i] = i
	}
	return data
}

func setupRandomAccess() ([]int, []int) {
	data := make([]int, dataSize)
	indices := make([]int, dataSize)
	for i := range data {
		data[i] = i
	}
	r := rand.New(rand.NewSource(42))
	for i := range indices {
		indices[i] = r.Intn(dataSize)
	}
	return data, indices
}

func setupNodes(n int) []*Node {
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = &Node{
			Value: i,
		}
		for j := range nodes[i].Data {
			nodes[i].Data[j] = byte(i + j)
		}
	}
	return nodes
}

const dataSize = 100_000_000

// BenchmarkComparisonLinearAccess 当数量足够庞大时，可以发现linear会比random快一些
func BenchmarkComparisonLinearAccess(b *testing.B) {
	b.Run("LinearAccess", func(b *testing.B) {
		data := setupLinearData()
		for b.Loop() {
			_ = SumLinear(data)
		}
	})
	b.Run("RandomAccess", func(b *testing.B) {
		data, indices := setupRandomAccess()
		for b.Loop() {
			_ = SumRandom(data, indices)
		}
	})
}

// BenchmarkComparisonPrefetch 当数量非常夸张时，可以发现prefetch效率远超过没有预加载的
func BenchmarkComparisonPrefetch(b *testing.B) {
	b.Run("WithoutPrefetch", func(b *testing.B) {
		nodes := setupNodes(dataSize)
		for b.Loop() {
			_ = ProcessWithoutPrefetch(nodes)
		}
	})
	b.Run("WithPrefetch", func(b *testing.B) {
		nodes := setupNodes(dataSize)
		for b.Loop() {
			_ = ProcessWithPrefetch(nodes)
		}
	})
}
