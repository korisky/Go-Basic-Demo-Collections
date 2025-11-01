package pref

import (
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
