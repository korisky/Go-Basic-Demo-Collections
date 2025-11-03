package pref

import "testing"

type Vec3Unaligned struct {
	X, Y, Z float32
}

func AddVectorsUnAligned(a, b, result []Vec3Unaligned) {
	for i := range len(a) {
		result[i].X = a[i].X + b[i].X
		result[i].Y = a[i].Y + b[i].Y
		result[i].Z = a[i].Z + b[i].Z
	}
}

// Vec3Aligned SIMD 在arm64使用的是NEON字符集128bit (不是cache-line的128Bytes)
// 不同的cpu使用的SIMD拓展不同, bit长度不同, 需要根据CPU进行优化
type Vec3Aligned struct {
	X, Y, Z float32
	_       float32 // padding for 16-byte alignment
}

func AddVectorsAligned(a, b, result []Vec3Aligned) {
	for i := range len(a) {
		result[i].X = a[i].X + b[i].X
		result[i].Y = a[i].Y + b[i].Y
		result[i].Z = a[i].Z + b[i].Z
	}
}

const vectorCount = 10000

func BenchmarkAddVectorsComparison(b *testing.B) {
	b.Run("SIMD-Unfriendly", func(b *testing.B) {
		a := make([]Vec3Unaligned, vectorCount)
		bb := make([]Vec3Unaligned, vectorCount)
		result := make([]Vec3Unaligned, vectorCount)
		for i := range a {
			a[i] = Vec3Unaligned{float32(i), float32(i * 2), float32(i * 3)}
			bb[i] = Vec3Unaligned{float32(i), float32(i), float32(i)}
		}
		b.ResetTimer()
		for range b.N {
			AddVectorsUnAligned(a, bb, result)
		}
	})
	b.Run("SIMD-Friendly", func(b *testing.B) {
		a := make([]Vec3Aligned, vectorCount)
		bb := make([]Vec3Aligned, vectorCount)
		result := make([]Vec3Aligned, vectorCount)
		for i := range a {
			a[i] = Vec3Aligned{X: float32(i), Y: float32(i * 2), Z: float32(i * 3)}
			bb[i] = Vec3Aligned{X: float32(i), Y: float32(i), Z: float32(i)}
		}
		b.ResetTimer()
		for range b.N {
			AddVectorsAligned(a, bb, result)
		}
	})
}
