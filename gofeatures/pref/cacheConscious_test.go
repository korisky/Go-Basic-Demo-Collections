package pref

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
	}
}
