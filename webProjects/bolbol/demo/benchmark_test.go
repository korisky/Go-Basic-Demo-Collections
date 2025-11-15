package demo

import (
	"context"
	repository2 "own/example/bolbol/demo/repository"
	serializer2 "own/example/bolbol/demo/serializer"
	"testing"
)

func BenchmarkArticle(b *testing.B) {
	ma := &mockArticle{items: map[int]repository2.Article{
		1: {
			ID:      1,
			Title:   "Title#1",
			Content: "content of the first article.",
		},
	}}
	a := serializer2.NewArticle(ma, 3)

	for b.Loop() {
		a.ByID(context.Background(), 10)
	}
}
