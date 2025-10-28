package pref

import (
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
	if capacity&(capacity-1) != 0 {
		panic("capacity must be power of 2")
	}
	return &RobinHoodMap{
		buckets:  make([]bucket, capacity),
		mask:     uint64(capacity - 1),
		capacity: capacity,
	}
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
)

func setupGoMap() map[uint64]uint64 {
	m := make(map[uint64]uint64, mapSize)
	for i := uint64(0); i < mapSize; i++ {
		m[i] = i * 2
	}
	return m
}

func setupRobinHoodMap() *RobinHoodMap {
	m := NewRobinHoodMap(2 * mapSize)
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

func BenchmarkComparisonOverMap(b *testing.B) {
	b.Run("GoMap-Lookup", func(b *testing.B) {
		m := setupGoMap()
		k := setupLookupKeys()
		b.ResetTimer()
		for range b.N {
			for _, key := range k {
				_, _ = m[key]
			}
		}
	})
	b.Run("RobinHood-Lookup", func(b *testing.B) {
		m := setupRobinHoodMap()
		k := setupLookupKeys()
		b.ResetTimer()
		for range b.N {
			for _, key := range k {
				_, _ = m.Get(key)
			}
		}
	})
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
			m := NewRobinHoodMap(2 * mapSize)
			for j := uint64(0); j < mapSize; j++ {
				m.Put(j, j*2)
			}
		}
	})
}
