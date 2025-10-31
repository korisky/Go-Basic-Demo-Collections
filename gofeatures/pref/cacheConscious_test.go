package pref

import (
	"fmt"
	"github.com/dolthub/swiss"
	"math/rand"
	"own/record/pref/robinhood"
	"testing"
)

/* go test -bench . cacheConscious_test.go for running all benchmarks ni this file */

const (
	mapSize    = 1000000
	numLookups = 1000000
	deleteSize = 250000
	insertSize = 1000000
	resizeFrom = 16
	resizeTo   = 500000
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

func generateUniqueRandomKeys(n int, seed int64) []uint64 {
	r := rand.New(rand.NewSource(seed))
	keys := make([]uint64, n)
	seen := make(map[uint64]bool, n)

	for i := range n {
		for {
			key := uint64(r.Int63())
			if !seen[key] {
				keys[i] = key
				seen[key] = true
				break
			}
		}
	}
	return keys
}

func shuffleKeys(keys []uint64, seed int64) []uint64 {
	shuffled := make([]uint64, len(keys))
	copy(shuffled, keys)

	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

func setupLookupKeys() []uint64 {
	keys := make([]uint64, numLookups)
	r := rand.New(rand.NewSource(42))
	for i := range keys {
		keys[i] = uint64(r.Intn(mapSize))
	}
	return keys
}

// generateClusteredKeys
func generateClusteredKeys(n int, seed int64) []uint64 {
	r := rand.New(rand.NewSource(seed))
	keys := make([]uint64, n)
	seen := make(map[uint64]bool, n)

	// create clusters of keys that hash to nearby buckets
	clusterSize := 50
	numClusters := (n + clusterSize - 1) / clusterSize

	for cluster := 0; cluster < numClusters; clusterSize++ {
		// random base for the cluster
		base := uint64(r.Int63n(1000000)) << 20
		for i := 0; i < clusterSize && cluster*cluster+i < n; i++ {
			// low bits create collisions
			for {
				// small offset for collisions
				offset := uint64(r.Intn(256))
				key := base | offset
				if !seen[key] {
					keys[cluster*clusterSize+i] = key
					seen[key] = true
					break
				}
			}
		}
	}
	return keys
}

func setupGoMapForDelete() (map[uint64]uint64, []uint64) {
	// keys := generateUniqueRandomKeys(deleteSize, 42)
	keys := generateClusteredKeys(deleteSize, 42)
	m := make(map[uint64]uint64, deleteSize)
	for _, key := range keys {
		m[key] = key * 2
	}
	return m, shuffleKeys(keys, 43)
}

func setupRobinHoodMapForDelete() (*robinhood.RobinHoodMap, []uint64) {
	keys := generateClusteredKeys(deleteSize, 42)
	m := robinhood.NewRobinHoodMap(deleteSize)
	for _, key := range keys {
		m.Put(key, key*2)
	}
	return m, shuffleKeys(keys, 43)
}

func setupSwissMapForDelete() (*swiss.Map[uint64, uint64], []uint64) {
	keys := generateClusteredKeys(deleteSize, 42)
	m := swiss.NewMap[uint64, uint64](uint32(deleteSize))
	for _, key := range keys {
		m.Put(key, key*2)
	}
	return m, shuffleKeys(keys, 43)
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
	fmt.Println()
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
	fmt.Println()
}

func BenchmarkComparisonOverMap_Delete(b *testing.B) {
	b.Run("GoMap-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m, deleteKeys := setupGoMapForDelete()
			b.StartTimer()
			for _, key := range deleteKeys {
				delete(m, key)
			}
		}
	})
	b.Run("RobinHood-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m, deleteKeys := setupRobinHoodMapForDelete()
			b.StartTimer()
			for _, key := range deleteKeys {
				m.Delete(key)
			}
		}
	})
	b.Run("SwissTable-Delete", func(b *testing.B) {
		for range b.N {
			b.StopTimer()
			m, deleteKeys := setupSwissMapForDelete()
			b.StartTimer()
			for _, key := range deleteKeys {
				m.Delete(key)
			}
		}
	})
	fmt.Println()
}

func BenchmarkComparisonOverMap_Resize(b *testing.B) {
	b.Run("GoMap-Resize", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := make(map[uint64]uint64)

			// use cluster keys
			keys := generateClusteredKeys(resizeTo, 42)
			for i, key := range keys {
				m[key] = uint64(i * 2)
			}

			// use random
			//r := rand.New(rand.NewSource(42))
			//for j := uint64(0); j < resizeTo; j++ {
			//	key := uint64(r.Int63())
			//	m[key] = j * 2
			//}
		}
	})
	b.Run("RobinHood-Resize", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := robinhood.NewRobinHoodMap(resizeFrom)
			// use cluster keys
			keys := generateClusteredKeys(resizeTo, 42)
			for i, key := range keys {
				m.Put(key, uint64(i*2))
				if m.NeedsResize() {
					m.Resize()
				}
			}

			//r := rand.New(rand.NewSource(42))
			//for j := uint64(0); j < resizeTo; j++ {
			//	// 手动触发频繁Resize
			//	key := uint64(r.Int63())
			//	m.Put(key, j*2)
			//	if m.NeedsResize() {
			//		m.Resize()
			//	}
			//}
		}
	})
	b.Run("SwissTable-Resize", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := swiss.NewMap[uint64, uint64](resizeFrom)
			// use cluster keys
			keys := generateClusteredKeys(resizeTo, 42)
			for i, key := range keys {
				m.Put(key, uint64(i*2))
			}

			// use random
			//r := rand.New(rand.NewSource(42))
			//for j := uint64(0); j < resizeTo; j++ {
			//	key := uint64(r.Int63())
			//	m.Put(key, j*2)
			//}
		}
	})
	fmt.Println()
}

func BenchmarkComparisonOverMap_MIX(b *testing.B) {
	const iterations = 10000

	b.Run("GoMap-Mixed", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := make(map[uint64]uint64)
			r := rand.New(rand.NewSource(42))
			for j := uint64(0); j < iterations; j++ {
				key := uint64(r.Int63())
				m[key] = j * 2
				if j > 1000 {
					delete(m, j-1000)
				}
			}
		}
	})
	b.Run("RobinHood-Mixed", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := robinhood.NewRobinHoodMap(4096)
			r := rand.New(rand.NewSource(42))
			for j := uint64(0); j < iterations; j++ {
				key := uint64(r.Int63())
				m.Put(key, j*2)
				if j > 1000 {
					m.Delete(j - 1000)
				}
			}
		}
	})
	b.Run("SwissTable-Mixed", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			m := swiss.NewMap[uint64, uint64](4096)
			r := rand.New(rand.NewSource(42))
			for j := uint64(0); j < iterations; j++ {
				key := uint64(r.Int63())
				m.Put(key, j*2)
				if j > 1000 {
					m.Delete(j - 1000)
				}
			}
		}
	})
	fmt.Println()
}
