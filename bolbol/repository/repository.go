package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("article not found")
)

type Article struct {
	ID      uint64
	Title   string
	Content string
}

// ArticleRepository is interface for manipulating Article's storage
type ArticleRepository interface {
	ByID(ctx context.Context, id int) (Article, error)
}

// SimpleSummaryArticle is a JSON format version or Article
type SimpleSummaryArticle struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	More    string `json:"more"`
}
