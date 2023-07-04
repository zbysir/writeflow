package repo

import (
	"context"
	"github.com/zbysir/writeflow/internal/model"
	"time"
)

type SearchArticleParams struct {
	Keyword     string    `json:"keyword"`
	Embedding   []float32 `json:"embedding"`
	BookIds     []int64   `json:"book_ids"`
	CategoryIds []int64   `json:"category_ids"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Offset      int       `json:"offset"`
	Limit       int       `json:"limit"`
}

type SearchFragmentParams struct {
	Keyword     string    `json:"keyword"`
	Embedding   []float32 `json:"embedding"`
	DocumentIds []int64   `json:"document_ids"`
	BookIds     []int64   `json:"book_ids"`
	MaxDistance float32   `json:"distance"`
	Offset      int       `json:"offset"`
	Limit       int       `json:"limit"`
}
type GetArticleListParams struct {
	BookIds     []int64 `json:"book_ids" form:"book_ids"`
	CategoryIds []int64 `json:"category_ids" form:"category_ids"`
	Offset      int     `json:"offset" form:"offset"`
	Limit       int     `json:"limit" form:"limit"`
}

type Document interface {
	GetDocumentList(ctx context.Context, p GetArticleListParams) (cs []model.Document, total int64, err error)
	// SaveDocument = Create + Update
	SaveDocument(ctx context.Context, content model.Document) (id int64, err error)
	DeleteDocument(ctx context.Context, id int64) (err error)
	SearchDocument(ctx context.Context, p SearchArticleParams) (cs []model.Document, total int64, err error)
	SearchFragment(ctx context.Context, p SearchFragmentParams) (cs []model.Fragment, err error)
}
