package robinhood

// ~24 bytes per bucket (increased from ~20 due to uint16 distance + padding)
type bucket struct {
	key      uint64
	value    uint64
	distance uint16
	occupied bool
}

type RobinHoodMap struct {
	buckets  []bucket
	mask     uint64
	size     int
	capacity int
}

func NewRobinHoodMap(capacity int) *RobinHoodMap {
	power2Cap := nextPowerOf2(int(float64(capacity) * 1.5))
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
	distance := uint16(0)

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
			b.distance = distance         // Update distance for newly inserted element
			distance = newBucket.distance // Continue with evicted element's distance
		}

		// move to nextSlot -> linear probing
		idx = (idx + 1) & m.mask
		distance++
	}
}

func (m *RobinHoodMap) Get(key uint64) (uint64, bool) {
	idx := key & m.mask
	distance := uint16(0)

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

// Delete RobinHoodMap's weakness, due to backward shifting
func (m *RobinHoodMap) Delete(key uint64) bool {
	idx := key & m.mask
	distance := uint16(0)

	for {
		b := &m.buckets[idx]

		if !b.occupied {
			return false
		}

		if b.key == key {
			b.occupied = false
			m.size--

			prevIdx := idx
			idx = (idx + 1) & m.mask

			// backward shift following entries
			for {
				curr := &m.buckets[idx]
				if !curr.occupied || curr.distance == 0 {
					// mark last shifted position as unoccupied
					m.buckets[prevIdx].occupied = false
					break
				}

				m.buckets[prevIdx] = *curr
				m.buckets[prevIdx].distance--
				curr.occupied = false

				prevIdx = idx
				idx = (idx + 1) & m.mask
			}
			return true
		}

		if b.distance < distance {
			return false
		}

		idx = (idx + 1) & m.mask
		distance++
	}
}

func (m *RobinHoodMap) NeedsResize() bool {
	return float64(m.size) > float64(m.capacity)*0.7
}

func (m *RobinHoodMap) Resize() {
	oldBuckets := m.buckets
	newCap := m.capacity * 2

	m.buckets = make([]bucket, newCap)
	m.mask = uint64(newCap - 1)
	m.capacity = newCap
	m.size = 0

	for _, b := range oldBuckets {
		if b.occupied {
			m.Put(b.key, b.value)
		}
	}
}

// Buckets are Diagnostic helpers for testing
func (m *RobinHoodMap) Buckets() []bucket {
	return m.buckets
}

func (m *RobinHoodMap) BucketAt(idx int) bucket {
	return m.buckets[idx]
}

func (b bucket) Occupied() bool {
	return b.occupied
}

func (b bucket) Distance() uint16 {
	return b.distance
}
