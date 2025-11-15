package pref

import "testing"

type Entity struct {
	ID       uint64
	X, Y, Z  float64
	Velocity float64
	Health   int
	Type     string
}

// World 是典型的 Arrays of Struct
type World struct {
	Entities []Entity
}

func (w *World) UpdatePositions(dt float64) {
	for i := range w.Entities {
		w.Entities[i].X += w.Entities[i].Velocity * dt
	}
}

// TheWorld 则是抽离的 Struct of Arrays
type TheWorld struct {
	IDs        []uint64
	Positions  [][3]float64
	Velocities []float64
	Healths    []int
	Types      []string
}

func (w *TheWorld) UpdatePositions(dt float64) {
	for i := range w.Positions {
		w.Positions[i][0] += w.Velocities[i] * dt
	}
}

// initWorlds 同时构建 World & TheWorld
func initWorlds(n int) (*World, *TheWorld) {
	// AoS
	world := &World{Entities: make([]Entity, n)}

	// SoA
	theWorld := &TheWorld{
		IDs:        make([]uint64, n),
		Positions:  make([][3]float64, n),
		Velocities: make([]float64, n),
		Healths:    make([]int, n),
		Types:      make([]string, n),
	}

	for i := range n {
		// AoS init
		world.Entities[i] = Entity{
			ID:       uint64(i),
			X:        float64(i),
			Y:        float64(i * 2),
			Z:        float64(i * 3),
			Velocity: 10.0,
			Health:   100,
			Type:     "enemy",
		}
		// SoA init
		theWorld.IDs[i] = uint64(i)
		theWorld.Positions[i] = [3]float64{float64(i), float64(i * 2), float64(i * 3)}
		theWorld.Velocities[i] = 10.0
		theWorld.Healths[i] = 100
		theWorld.Types[i] = "enemy"
	}

	return world, theWorld
}

const (
	numEntities = 1000000
	deltaTime   = 0.016 // 60 FPS
)

// BenchmarkComparisonSoAAndAoS 可以看出, SoA什么情况下都会比AoS要更快
// 随着Array的数量越大, 或者struct越复杂, 性能提升越明显 (当前struct+100w数量, 性能提升有35%)
func BenchmarkComparisonSoAAndAoS(b *testing.B) {
	worlds, theWorld := initWorlds(numEntities)
	b.Run("AoS", func(b *testing.B) {
		for b.Loop() {
			worlds.UpdatePositions(deltaTime)
		}
	})
	b.Run("SoA", func(b *testing.B) {
		for b.Loop() {
			theWorld.UpdatePositions(deltaTime)
		}
	})
}
