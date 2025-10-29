package pref

import (
	"github.com/dolthub/swiss"
	"math/rand"
	"own/record/pref/robinhood"
	"testing"
)

const (
	mapSize    = 1000000
	numLookups = 1000000
	deleteSize = 500000
	insertSize = 500000
	resizeFrom = 16
	resizeTo   = 100000
)

func setupGoMap() map[uint64]uint64 {
	m := make(map[uint64]uint64, mapSize)
	for i := uint64(0); i < mapSize; i++ {
		m[i] = i * 2
	}
	return m
}

func setupRobinHoodMap() *robinhood.RobinHoodMap {
	m := robinhood.NewRobinHoodMap(mapSize)
	for i := uint64(0); i < mapSize; i++ {
		m.Put(i, i*2)
	}
	return m
}

func setupSwissMap() *swiss.Map[uint64, uint64] {
	m := swiss.NewMap[uint64, uint64](uint32(mapSize))
	for i := uint64(0); i < mapSize; i++ {
		m.Put(i, i*2)
	}
	return m
}

func setupLookupKeys() []uint64 {
	keys := make([]uint64, numLookups)
	r := rand.New(rand.NewSource(42))
	for i := range keys {
		keys[i] = uint64(r.Intn(mapSize))
	}
	return keys
}

func setupGoMapForDelete() map[uint64]uint64 {
	m := make(map[uint64]uint64, deleteSize)
	for i := uint64(0); i < deleteSize; i++ {
		m[i] = i * 2
	}
	return m
}

func setupRobinHoodMapForDelete() *robinhood.RobinHoodMap {
	m := robinhood.NewRobinHoodMap(deleteSize)
	for i := uint64(0); i < deleteSize; i++ {
		m.Put(i, i*2)
	}
	return m
}

func setupSwissMapForDelete() *swiss.Map[uint64, uint64] {
	m := swiss.NewMap[uint64, uint64](uint32(deleteSize))
	for i := uint64(0); i < deleteSize; i++ {
		m.Put(i, i*2)
	}
	return m
}

// BenchmarkComparisonOverMap 这里可以看出 RobinHood Map 对比 Map 实现性能非常夸张
// 但为何 RobinHood Map 不是 Map 的默认实现呢？
// 1) Key的限制, RobinHoodMap 限制了key一定是无符号整型
// 2) RobinHoodMap 基本不支持delete(非常重的操作, resize同样）， 而下面的unitTest更多的insert & randomRead, 都RobinHoodMap的强项
// 3) RobinHoodMap 对于loadFactor非常敏感, 过大(>90%)就能使得性能急剧下降, 所以已知大小的情况一般空间换时间要2x
// 4) RobinHoodMap 不支持并发, 要使其支持并发也会增加很多overhead, 而go中原map实现通过shard分片, 使得并发成为可能
// Rust语言的Map使用 RobinHoodMap ？
// 2019之前Rust使用的就是 RobinHoodMap 实现的map, 随后改成Google的 SwissTable (SIMD lookup + easier delete/resize)
func BenchmarkComparisonOverMap_LookUp(b *testing.B) {

	goMap := setupGoMap()
	robinMap := setupRobinHoodMap()
	swissMap := setupSwissMap()
	keys := setupLookupKeys()

	b.Run("GoMap-Lookup", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			for _, key := range keys {
				_, _ = goMap[key]
			}
		}
	})
	b.Run("RobinHood-Lookup", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			for _, key := range keys {
				_, _ = robinMap.Get(key)
			}
		}
	})
	b.Run("SwissTable-Lookup", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			for _, key := range keys {
				_, _ = swissMap.Get(key)
			}
		}
	})
}

func BenchmarkComparisonOverMap_Insert(b *testing.B) {
	b.Run("GoMap-Insert", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := make(map[uint64]uint64, mapSize)
			for j := uint64(0); j < mapSize; j++ {
				m[j] = j * 2
			}
		}
	})
	b.Run("RobinHood-Insert", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := robinhood.NewRobinHoodMap(insertSize)
			for j := uint64(0); j < mapSize; j++ {
				m.Put(j, j*2)
			}
		}
	})
	b.Run("SwissTable-Insert", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := swiss.NewMap[uint64, uint64](uint32(mapSize))
			for j := uint64(0); j < mapSize; j++ {
				m.Put(j, j*2)
			}
		}
	})
}

func BenchmarkComparisonOverMap_Delete(b *testing.B) {
	b.Run("GoMap-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m := setupGoMapForDelete()
			b.StartTimer()
			for j := uint64(0); j < deleteSize; j++ {
				delete(m, j)
			}
		}
	})
	b.Run("RobinHood-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m := setupRobinHoodMapForDelete()
			b.StartTimer()
			for j := uint64(0); j < deleteSize; j++ {
				m.Delete(j)
			}
		}
	})
	b.Run("SwissTable-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m := setupSwissMapForDelete()
			b.StartTimer()
			for j := uint64(0); j < deleteSize; j++ {
				m.Delete(j)
			}
		}
	})
}
