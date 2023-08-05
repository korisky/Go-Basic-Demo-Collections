package storage

import (
	"context"
	"own/example/bolbol/entity"
	"sync"
)

type memoryWithChannel struct {
	storage *sync.Map
	size    int
}

func NewMemoryWithChannel(size int) Storage {
	return &memoryWithChannel{
		storage: new(sync.Map),
		size:    size,
	}
}

func (m *memoryWithChannel) Push(ctx context.Context, clientID int, notification entity.Notification) error {
	c := m.get(clientID)
	if len(c) == m.size {
		<-c // remove the latest item and requires a garbage collection to be deleted from the heap
	}
	c <- notification
	return nil
}

func (m *memoryWithChannel) Count(ctx context.Context, clientID int) (int, error) {
	c := m.get(clientID)
	return len(c), nil
}

func (m *memoryWithChannel) Pop(ctx context.Context, clientID int) (entity.Notification, error) {
	c := m.get(clientID)
	select {
	case item := <-c:
		return item, nil
	default:
		return nil, ErrEmpty
	}
}

func (m *memoryWithChannel) PopAll(ctx context.Context, clientID int) ([]entity.Notification, error) {
	c := m.get(clientID)
	l := len(c)
	items := make([]entity.Notification, l)
	for i := 0; i < l; i++ {
		items[i] = <-c
	}
	return items, nil
}

func (m *memoryWithChannel) get(clientID int) chan entity.Notification {
	cInf, _ := m.storage.LoadOrStore(clientID, make(chan entity.Notification, m.size))
	return cInf.(chan entity.Notification)
}
