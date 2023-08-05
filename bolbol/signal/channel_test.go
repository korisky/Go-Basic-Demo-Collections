package signal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChannel(t *testing.T) {
	s := NewSignal()
	ac, _, err := s.Subscribe("a")
	assert.NoError(t, err)

	ac2, _, err := s.Subscribe("a")
	assert.NoError(t, err)

	_ = s.Publish("a")

	select {
	case <-ac:
	default:
		t.Fatal("didn't receive the signal")
	}
	select {
	case <-ac2:
	default:
		t.Fatal("didn't receive the signal")
	}

	err = s.Publish("b")
	assert.ErrorIs(t, err, ErrEmpty)

	_, cancel, err := s.Subscribe("c")
	assert.NoError(t, err)
	cancel()
	err = s.Publish("c")
	assert.ErrorIs(t, err, ErrEmpty)
}
