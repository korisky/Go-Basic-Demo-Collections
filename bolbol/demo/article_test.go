package demo

import (
	"context"
	"github.com/stretchr/testify/assert"
	repository2 "own/example/bolbol/demo/repository"
	serializer2 "own/example/bolbol/demo/serializer"
	"testing"
)

type mockArticle struct {
	items map[int]repository2.Article
}

func (m *mockArticle) ByID(ctx context.Context, id int) (repository2.Article, error) {
	val, has := m.items[id]
	if !has {
		return repository2.Article{}, repository2.ErrNotFound
	}
	return val, nil
}

func TestArticle_ByID(t *testing.T) {

	ma := &mockArticle{items: map[int]repository2.Article{
		1: {
			ID:      1,
			Title:   "Title#1",
			Content: "content of the first article.",
		},
	}}
	a := serializer2.NewArticle(ma, 3)

	_, err := a.ByID(context.Background(), 10)
	assert.ErrorIs(t, repository2.ErrNotFound, err)

	item, err := a.ByID(context.Background(), 1)
	assert.Equal(t, "https://site.com/a/1", item.More)
	assert.Equal(t, uint64(1), item.ID)
	assert.Equal(t, "content of the", item.Summary)
}
