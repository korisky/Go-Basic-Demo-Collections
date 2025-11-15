package pref

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	ActionView uint8 = iota
	ActionClick
	ActionPurchase
	ActionShare
)

var totalRevenue float64

func updateRevene(val float64) {
	totalRevenue += val
}

type Event struct {
	Timestamp time.Time
	UserId    uint64
	Action    string
	Value     float64
	Tags      map[string]string
}

func ProcessEvents(events []Event) {
	for _, event := range events {
		if event.Action == "purchase" {
			updateRevene(event.Value)
		}
	}
}

// EventBatch Data-Oriented design, SoA
type EventBatch struct {
	Timestamps []int64
	UserIds    []uint64
	Actions    []uint8
	Values     []float64

	TagIndices []uint32
	TagKeys    []string
	TagValues  []string
}

func ProcessEventsBatch(batch *EventBatch) {
	for i, action := range batch.Actions {
		if action == ActionPurchase {
			updateRevene(batch.Values[i])
		}
	}
}

func genEvents(n int) []Event {
	r := rand.New(rand.NewSource(42))
	events := make([]Event, n)
	actions := []string{"view", "click", "purchase", "share"}
	for i := range n {
		events[i] = Event{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			UserId:    r.Uint64(),
			Action:    actions[r.Intn(len(actions))],
			Value:     r.Float64() * 1000,
			Tags: map[string]string{
				"source":   "web",
				"campaign": "summer_sale",
			},
		}
	}
	return events
}

func genEventBatch(n int) *EventBatch {
	r := rand.New(rand.NewSource(42))
	batch := &EventBatch{
		Timestamps: make([]int64, n),
		UserIds:    make([]uint64, n),
		Actions:    make([]uint8, n),
		Values:     make([]float64, n),
	}
	for i := range n {
		batch.Timestamps[i] = time.Now().Add(time.Duration(i) * time.Second).Unix()
		batch.UserIds[i] = r.Uint64()
		batch.Actions[i] = uint8(r.Intn(4))
		batch.Values[i] = r.Float64() * 1000
	}
	return batch
}

func BenchmarkAnalyticsPipeline(b *testing.B) {
	sizes := []int{1000, 10000, 100000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("OOP_%dk", size/1000), func(b *testing.B) {
			events := genEvents(size)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				ProcessEvents(events)
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*size), "ns/events")
		})
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("DOD_%dk", size/1000), func(b *testing.B) {
			events := genEventBatch(size)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				ProcessEventsBatch(events)
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*size), "ns/events")
		})
	}
}
