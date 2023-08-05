package bolbol

import (
	"context"
	"github.com/stretchr/testify/assert"
	"own/example/bolbol/repository"
	"own/example/bolbol/serializer"
	"testing"
)

type mockArticle struct {
	items map[int]repository.Article
}

func (m *mockArticle) ByID(ctx context.Context, id int) (repository.Article, error) {
	val, has := m.items[id]
	if !has {
		return repository.Article{}, repository.ErrNotFound
	}
	return val, nil
}

func TestArticle_ByID(t *testing.T) {

	ma := &mockArticle{items: map[int]repository.Article{
		1: {
			ID:      1,
			Title:   "Title#1",
			Content: "content of the first article.",
		},
	}}
	a := serializer.NewArticle(ma, 3)

	_, err := a.ByID(context.Background(), 10)
	assert.ErrorIs(t, repository.ErrNotFound, err)

	item, err := a.ByID(context.Background(), 1)
	assert.Equal(t, "https://site.com/a/1", item.More)
	assert.Equal(t, uint64(1), item.ID)
	assert.Equal(t, "content of the", item.Summary)
}
