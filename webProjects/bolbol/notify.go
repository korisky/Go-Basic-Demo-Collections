package bolbol

import (
	"context"
	"fmt"
	"own/example/bolbol/entity"
	"strconv"
)

func (b *Bolbol) Notify(ctx context.Context, userID int, notification entity.Notification) error {
	if err := b.Storage.Push(ctx, userID, notification); err != nil {
		return fmt.Errorf("error while trying to push the new notification: %w", err)
	}
	_ = b.Signal.Publish("user#" + strconv.Itoa(userID))
	return nil
}
