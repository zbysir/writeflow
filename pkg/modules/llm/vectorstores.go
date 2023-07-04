package llm

import "context"

type Fragment struct {
	Body   string                 `json:"body"`
	Meta   map[string]interface{} `json:"meta"`
	Vector []float32              `json:"vector"`
}

type SimilaritySearchParams struct {
	Vector []float32
	Number int
}

type VectorStoreFactory interface {
	NewVectorStore(ctx context.Context, config map[string]interface{}) (VectorStore, error)
}

type VectorStore interface {
	SimilaritySearch(ctx context.Context, p SimilaritySearchParams) (fs []Fragment, err error)
}
