package bolbol

import (
	"context"
	"own/example/bolbol/repository"
	"own/example/bolbol/serializer"
	"testing"
)

func BenchmarkArticle(b *testing.B) {
	ma := &mockArticle{items: map[int]repository.Article{
		1: {
			ID:      1,
			Title:   "Title#1",
			Content: "content of the first article.",
		},
	}}
	a := serializer.NewArticle(ma, 3)

	for i := 0; i < b.N; i++ {
		a.ByID(context.Background(), 10)
	}
}
