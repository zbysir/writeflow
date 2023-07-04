package repo

import (
	"context"
	"fmt"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/modules/llm"
	"testing"
	"xorm.io/xorm"
)

func setup(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=123456 dbname=postgres sslmode=disable")
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		return
	}
	engine.ShowSQL(true)
	engine.DropTables(model.Document{}, model.Fragment{})

	_, err = engine.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Exec("CREATE TABLE IF NOT EXISTS \"public\".\"fragment\"\n(\n    \"id\"\n    BIGSERIAL\n    PRIMARY\n    KEY\n    NOT\n    NULL,\n    \"article_id\"\n    BIGINT\n    NULL,\n    \"body\"\n    VARCHAR(    255) NULL,\n    \"start_index\" INTEGER NULL, \"end_index\" INTEGER NULL, \"vector\" vector(3) , \"created_at\" TIMESTAMP NULL, \"updated_at\" TIMESTAMP NULL);")
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Exec("CREATE TABLE IF NOT EXISTS \"public\".\"article\" (\"id\" BIGSERIAL PRIMARY KEY  NOT NULL, \"book_id\" BIGINT NULL, \"category_id\" BIGINT NULL, \"title\" VARCHAR(255) NULL, \"body\" VARCHAR(255) NULL, \"embedding_status\" VARCHAR(255) NULL, \"created_at\" TIMESTAMP NULL, \"updated_at\" TIMESTAMP NULL);")
	if err != nil {
		t.Fatal(err)
	}
	//err = engine.Sync(Document{})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = engine.Sync(Fragment{})
	//if err != nil {
	//	t.Fatal(err)
	//}
}

type mockEmbeddinger struct {
}

func (m mockEmbeddinger) Embedding(input []string) ([][]float32, error) {
	r := [][]float32{}
	for _ = range input {
		r = append(r, []float32{1, 1, 3})
	}
	return r, nil
}

func TestSaveArticle(t *testing.T) {
	setup(t)

	s, err := NewPGStorage(DbConfig{
		Debug:    true,
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
		DBName:   "postgres",
	}, llm.NewMarkDoneSplit(1024), mockEmbeddinger{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	id, err := s.SaveDocument(ctx, model.Document{
		Id: 0,
		Source: model.DocumentSource{
			Body: "hello world",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.CreateEmbedding(ctx, []int64{id})
	if err != nil {
		t.Fatal(err)
	}

	a, err := s.GetDocument(ctx, id, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id: %v", id)
	t.Logf("article: %v", a)

	fs, err := s.SearchFragment(ctx, SearchFragmentParams{
		Keyword:     "",
		Embedding:   []float32{1, 1, 3},
		DocumentIds: nil,
		MaxDistance: 0,
		Offset:      0,
		Limit:       0,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Fragment: %v", fs)
}
