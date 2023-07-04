package repo

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/pkg/modules/llm"
	"strings"
	"time"
	"xorm.io/xorm"
)

import (
	_ "github.com/lib/pq"
)

type PGStorage struct {
	orm         *xorm.Engine
	spliter     llm.Spliter
	embeddinger llm.Embeddinger
}

var _ Document = (*PGStorage)(nil)

func (p *PGStorage) GetDocumentList(ctx context.Context, pr GetArticleListParams) (cs []model.Document, total int64, err error) {
	return nil, 0, nil
}

func (p *PGStorage) GetDocument(ctx context.Context, id int64, withFragment bool) (cs model.Document, err error) {
	exist, err := p.orm.Context(ctx).Where("id=?", id).
		Get(&cs)
	if err != nil {
		return model.Document{}, fmt.Errorf("orm.Get Document error: %w", err)
	}
	if !exist {
		return model.Document{}, nil
	}

	if withFragment {
		var fs []model.Fragment
		err = p.orm.Context(ctx).Where("document_id=?", id).Find(&fs)
		if err != nil {
			return model.Document{}, err
		}

		cs.Fragments = fs
	}

	return
}

func (p *PGStorage) SaveDocument(ctx context.Context, content model.Document) (id int64, err error) {
	if content.Id == 0 {
		_, err = p.orm.Context(ctx).Insert(&content)
		if err != nil {
			return
		}
		id = content.Id
	} else {
		_, err = p.orm.Context(ctx).Where("id=?", content.Id).Update(&content)
		if err != nil {
			return
		}
	}

	return
}

func (p *PGStorage) DeleteDocument(ctx context.Context, id int64) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p *PGStorage) AsyncEmbedding(ctx context.Context, articleIds []int64) (err error) {
	embeddingArticle := func(a *model.Document) error {
		ss := p.spliter.Split(a.Source.Body)
		fs := make([]model.Fragment, 0, len(ss))
		ver, err := p.embeddinger.Embedding(ss)
		if err != nil {
			return err
		}

		startIndex := 0
		for i, s := range ss {
			end := startIndex + len([]rune(s))
			fs = append(fs, model.Fragment{
				Id:         0,
				DocumentId: a.Id,
				Body:       s,
				StartIndex: startIndex,
				EndIndex:   end,
				Vector:     ver[i],
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
			})
			startIndex = end
		}

		_, err = p.orm.Context(ctx).Omit("similarity").Insert(&fs)
		if err != nil {
			return fmt.Errorf("orm.Insert Fragment error: %w", err)
		}

		return nil
	}

	for _, v := range articleIds {
		a, err := p.GetDocument(ctx, v, true)
		if err != nil {
			return err
		}

		err = embeddingArticle(&a)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PGStorage) SearchDocument(ctx context.Context, keyword SearchArticleParams) (cs []model.Document, total int64, err error) {
	session := p.orm.Context(ctx)
	if len(keyword.Embedding) != 0 {
		var fs []model.Fragment
		err = p.orm.Context(ctx).Where("vector<-> ?", keyword.Embedding).Find(&fs)
		if err != nil {
			return
		}
		session = session.Where("in in ?", lo.Map(fs, func(item model.Fragment, index int) int64 {
			return item.DocumentId
		}))
	}

	session = session.Limit(keyword.Limit, keyword.Offset)

	total, err = session.FindAndCount(&cs)
	if err != nil {
		return
	}

	return
}

func floatSliceToString(s []float32) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, item := range s {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%v", item))
	}
	sb.WriteString("]")
	return sb.String()
}

func (p *PGStorage) SearchFragment(ctx context.Context, params SearchFragmentParams) (fs []model.Fragment, err error) {
	session := p.orm.Context(ctx)
	fs = []model.Fragment{}
	if params.Limit != 0 {
		session = session.Limit(params.Limit, params.Offset)
	}
	if len(params.Embedding) != 0 {
		if params.MaxDistance != 0 {
			session = session.Where("vector<-> ? <= ?", floatSliceToString(params.Embedding), params.MaxDistance)
		}
		session = session.Select(fmt.Sprintf("id, body, 1 - (vector<=> '%v') AS similarity", floatSliceToString(params.Embedding)))
		session = session.OrderBy("similarity desc")
	}
	err = session.Find(&fs)
	if err != nil {
		return
	}

	return
}

func (p *PGStorage) GetFragmentDistance(ctx context.Context, embedding []float32, maxDistance float64) (fs []model.Fragment, err error) {
	session := p.orm.Context(ctx)
	fs = []model.Fragment{}
	if maxDistance != 0 {
		session = session.Where("vector<=> ? <= ?", floatSliceToString(embedding), maxDistance)
	}

	session = session.Select(fmt.Sprintf("id, body, 1 - (vector<=> '%v') AS similarity", floatSliceToString(embedding)))

	err = session.OrderBy("similarity desc").Find(&fs)
	if err != nil {
		return
	}

	return
}

type NewPGStorageParams struct {
	DbConfig    DbConfig
	spliter     llm.Spliter
	embeddinger llm.Embeddinger
}
type DbConfig struct {
	Debug    bool
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewPGStorage(params DbConfig, spliter llm.Spliter, embeddinger llm.Embeddinger) (s *PGStorage, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", params.Host, params.Port, params.User, params.Password, params.DBName)
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("xorm.NewEngine error: %w", err)
	}
	go func() {
		_, err = engine.Exec("CREATE EXTENSION IF NOT EXISTS vector")
		if err != nil {
			log.Errorf("[Postgres] enable vector extension error: %v", err)
		}
	}()

	engine.ShowSQL(!params.Debug)
	return &PGStorage{
		orm:         engine,
		spliter:     spliter,
		embeddinger: embeddinger,
	}, nil
}
