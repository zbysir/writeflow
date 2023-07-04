package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/modules/llm"
	"os"
	"strings"
)

func main() {
	cli := openai.NewClient(os.Getenv("APIKEY"))

	s, err := repo.NewPGStorage(repo.DbConfig{
		Debug:    true,
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
		DBName:   "writeflow",
	}, llm.NewMarkDoneSplit(1024*2), llm.NewOpenAIEmbedding(cli))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	r, err := cli.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{"计算余弦相似度的操作符是什么"},
		Model: openai.AdaEmbeddingV2,
		User:  "",
	})
	if err != nil {
		panic(err)
	}

	fs, err := s.GetFragmentDistance(ctx, r.Data[0].Embedding, 1000)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		fmt.Printf("id: %v, similarity:%+v, %v\n\n", f.Id, f.Similarity, strings.ReplaceAll(f.Body, "\n", " "))
	}
}
