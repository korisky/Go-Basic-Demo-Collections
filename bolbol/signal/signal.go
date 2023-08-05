package signal

import "errors"

var (
	ErrEmpty = errors.New("no topic found")
)

type Signal interface {
	Subscribe(id string) (<-chan struct{}, func(), error)
	Publish(id string) error
}
