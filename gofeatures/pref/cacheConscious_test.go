package pref

import (
	"github.com/dolthub/swiss"
	"math/rand"
	"testing"
)

// ~20 bytes per bucket
type bucket struct {
	key      uint64
	value    uint64
	distance uint8
	occupied bool
}

type RobinHoodMap struct {
	buckets  []bucket
	mask     uint64
	size     int
	capacity int
}

func NewRobinHoodMap(capacity int) *RobinHoodMap {
	power2Cap := nextPowerOf2(capacity * 2)
	return &RobinHoodMap{
		buckets:  make([]bucket, power2Cap),
		mask:     uint64(power2Cap - 1),
		capacity: power2Cap,
	}
}

func nextPowerOf2(n int) int {
	if n <= 0 {
		return 1
	}
	if n&(n-1) == 0 {
		return n
	}
	power := 1
	for power < n {
		power <<= 1
	}
	return power
}

func (m *RobinHoodMap) Put(key, value uint64) {
	idx := key & m.mask
	distance := uint8(0)

	newBucket := bucket{
		key:      key,
		value:    value,
		distance: 0,
		occupied: true,
	}

	// liner probing
	for {
		b := &m.buckets[idx]

		// empty slot -> insert directly
		if !b.occupied {
			*b = newBucket
			b.distance = distance
			m.size++
			return
		}

		// key exists -> update val
		if b.key == key {
			b.value = value
			return
		}

		// Robin Hood: if current bucket's distance < new val's distance
		// swap it
		if b.distance < distance {
			newBucket, *b = *b, newBucket // golang can do this
			distance = b.distance
		}

		// move to nextSlot -> linear probing
		idx = (idx + 1) & m.mask
		distance++
	}
}

func (m *RobinHoodMap) Get(key uint64) (uint64, bool) {
	idx := key & m.mask
	distance := uint8(0)

	// linear probing
	for {
		b := &m.buckets[idx]

		// empty
		if !b.occupied {
			return 0, false
		}

		// found
		if b.key == key {
			return b.value, true
		}

		// Robin Hood Optimization: if current dis < search dis, must not exist
		if b.distance < distance {
			return 0, false
		}

		idx = (idx + 1) & m.mask
		distance++
	}
}

const (
	mapSize    = 1000000
	numLookups = 1000000
	deleteSize = 500000
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

func setupRobinHoodMap() *RobinHoodMap {
	m := NewRobinHoodMap(mapSize)
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
			m := NewRobinHoodMap(mapSize)
			for j := uint64(0); j < mapSize; j++ {
				m.Put(j, j*2)
			}
		}
	})
}
