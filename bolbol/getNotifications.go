package bolbol

import (
	"context"
	"fmt"
	"own/example/bolbol/entity"
	"strconv"
)

func (b *Bolbol) GetNotifications(ctx context.Context, clientID int) ([]entity.Notification, error) {
	c, err := b.Storage.Count(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("error while counting user's notifications: %w", err)
	}
	if c > 0 {
		return b.Storage.PopAll(ctx, clientID)
	}

	ch, cancel, err := b.Signal.Subscribe("user#" + strconv.Itoa(clientID))
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf("error while trying to listen on notification topic: %w", err)
	}

	ctx, ctxCancel := context.WithTimeout(ctx, b.defaultTimeout)
	defer ctxCancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return b.Storage.PopAll(ctx, clientID)
	}
}
