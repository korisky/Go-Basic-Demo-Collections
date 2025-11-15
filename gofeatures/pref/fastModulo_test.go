package pref

import "testing"

func modulo(key, capacity uint64) uint64 {
	return key % capacity
}

// fastModulo 对于binary的操作而言, modulo = 仅保留低位的bits
// 例如: number = (quotient × 2^k) + remainder
// 42 % 16 = 10   42 = 0b00101010, 16 = 0b00010000, 16-1=15=0b00001111
// 10 = 0b00001010 -> 42 & 15
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
		for b.Loop() {
			_ = modulo(key, capacity)
		}
	})
	b.Run("Modulo-VaryingKeys", func(b *testing.B) {
		i := 0
		for b.Loop() {
			_ = modulo(keys[i%len(keys)], capacity)
			i++
		}
	})
	b.Run("FastModulo-SingleKey", func(b *testing.B) {
		for b.Loop() {
			_ = fastModulo(key, capacity)
		}
	})
	b.Run("FastModulo-VaryingKeys", func(b *testing.B) {
		i := 0
		for b.Loop() {
			_ = fastModulo(keys[i%len(keys)], capacity)
			i++
		}
	})
}
