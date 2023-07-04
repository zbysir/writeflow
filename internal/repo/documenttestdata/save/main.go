package main

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/modules/llm"
	"os"

	"time"
	"xorm.io/xorm"
)

func setupDb() error {
	engine, err := xorm.NewEngine("postgres", "host=localhost port=5432 user=postgres password=123456 dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	engine.ShowSQL(true)
	_, _ = engine.Exec("CREATE DATABASE writeflow")
	engine, err = xorm.NewEngine("postgres", "host=localhost port=5432 user=postgres password=123456 dbname=writeflow sslmode=disable")
	if err != nil {
		return err
	}

	_, err = engine.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		return err
	}
	_, err = engine.Exec("CREATE TABLE IF NOT EXISTS \"public\".\"fragment\"\n(\n    \"id\"\n    BIGSERIAL\n    PRIMARY\n    KEY\n    NOT\n    NULL,\n    \"article_id\"\n    BIGINT\n    NULL,\n    \"body\"\n    VARCHAR(    255) NULL,\n    \"start_index\" INTEGER NULL, \"end_index\" INTEGER NULL, \"vector\" vector(3) , \"created_at\" TIMESTAMP NULL, \"updated_at\" TIMESTAMP NULL);")
	if err != nil {
		return err
	}
	_, err = engine.Exec(`
CREATE TABLE IF NOT EXISTS "public"."article"
(
    "id"               BIGSERIAL PRIMARY KEY NOT NULL,
    "book_id"          BIGINT                NULL,
    "category_id"      BIGINT                NULL,
    "title"            VARCHAR(255)          NULL,
    "source"           jsonb                 NULL,
    "embedding_status" VARCHAR(255)          NULL,
    "created_at"       TIMESTAMP             NULL,
    "updated_at"       TIMESTAMP             NULL
);`)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	body, err := os.ReadFile("/Users/bysir/goproj/bysir/writeflow/internal/repo/documenttestdata/about_bysir.md")
	if err != nil {
		panic(err)
	}

	err = setupDb()
	if err != nil {
		panic(err)
	}

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
	id, err := s.SaveDocument(ctx, model.Document{
		Id:     0,
		BookId: 0,
		Title:  "",
		Source: model.DocumentSource{
			Location: "",
			Path:     "",
			Body:     string(body),
		},
		EmbeddingStatus: "",
		Fragments:       nil,
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	})
	if err != nil {
		panic(err)
	}

	err = s.AsyncEmbedding(ctx, []int64{id})
	if err != nil {
		panic(err)
	}
}
