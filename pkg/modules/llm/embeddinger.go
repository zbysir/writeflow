package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Embeddinger interface {
	Embedding(input []string) ([][]float32, error)
}

type OpenAIEmbedding struct {
	cli *openai.Client
}

func NewOpenAIEmbedding(cli *openai.Client) *OpenAIEmbedding {
	return &OpenAIEmbedding{cli: cli}
}

func (o *OpenAIEmbedding) Embedding(input []string) ([][]float32, error) {
	rsp, err := o.cli.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: input,
		Model: openai.AdaEmbeddingV2,
		User:  "",
	})
	if err != nil {
		return nil, err
	}

	var ver [][]float32
	for _, v := range rsp.Data {
		ver = append(ver, v.Embedding)
	}

	return ver, nil
}
