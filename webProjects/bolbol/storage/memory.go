package storage

import (
	"context"
	"own/example/bolbol/entity"
	"sync"
)

type userStorage struct {
	mu            *sync.Mutex
	notifications []entity.Notification
}

type memoryWithList struct {
	size    int
	storage *sync.Map
}

func NewMemoryWithList(size int) Storage {
	return &memoryWithList{
		size:    size,
		storage: new(sync.Map),
	}
}

func (m *memoryWithList) Push(ctx context.Context, clientID int, notification entity.Notification) error {
	item := m.get(ctx, clientID)
	item.mu.Lock()
	defer item.mu.Unlock()
	if len(item.notifications) == m.size {
		item.notifications = item.notifications[1:]
	}
	item.notifications = append(item.notifications, notification)
	return nil
}

func (m *memoryWithList) Count(ctx context.Context, clientID int) (int, error) {
	item := m.get(ctx, clientID)
	return len(item.notifications), nil // len is a thread-safe function
}

func (m *memoryWithList) Pop(ctx context.Context, clientID int) (entity.Notification, error) {
	item := m.get(ctx, clientID)
	if len(item.notifications) == 0 {
		return nil, ErrEmpty
	}
	item.mu.Lock()
	defer item.mu.Unlock()

	notification := item.notifications[0]
	item.notifications = item.notifications[1:]
	return notification, nil
}

func (m *memoryWithList) PopAll(ctx context.Context, clientID int) ([]entity.Notification, error) {
	item := m.get(ctx, clientID)
	if len(item.notifications) == 0 {
		return nil, ErrEmpty
	}
	item.mu.Lock()
	defer item.mu.Unlock()
	defer func() {
		item.notifications = make([]entity.Notification, 0, m.size)
	}()

	return item.notifications, nil
}
func (m *memoryWithList) get(ctx context.Context, clientID int) *userStorage {
	item, _ := m.storage.LoadOrStore(clientID, &userStorage{mu: new(sync.Mutex), notifications: make([]entity.Notification, 0, m.size)})
	return item.(*userStorage)
}
