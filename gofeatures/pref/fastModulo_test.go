package pref

import "testing"

func modulo(key, capacity uint64) uint64 {
	return key % capacity
}

func fastModulo(key, capacity uint64) uint64 {
	if capacity&(capacity-1) != 0 {
		panic("capacity must be power of 2 for fast modulo")
	}
	mask := capacity - 1
	return key & mask
}

func BenchmarkModuloComparison(b *testing.B) {

	capacity := uint64(4096)
	key := uint64(123456789123456789)

	keys := make([]uint64, 1000)
	for i := range keys {
		keys[i] = uint64(i * 12345)
	}

	b.Run("Modulo-SingleKey", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_ = modulo(key, capacity)
		}
	})
	b.Run("FastModulo-SingleKey", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_ = fastModulo(key, capacity)
		}
	})
	b.Run("Modulo-VaryingKeys", func(b *testing.B) {
		b.ResetTimer()
		for i := range b.N {
			_ = modulo(keys[i%len(keys)], capacity)
		}
	})
	b.Run("FastModulo-VaryingKeys", func(b *testing.B) {
		b.ResetTimer()
		for i := range b.N {
			_ = fastModulo(keys[i%len(keys)], capacity)
		}
	})
}
