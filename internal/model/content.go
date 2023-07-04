package model

import "time"

type Book struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//type BookCategory struct {
//	Id        int64     `json:"id"`
//	BookId    int64     `json:"book_id"`
//	ParentId  int64     `json:"parent_id"`
//	Name      string    `json:"name"`
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}

// book -> document -> fragment

type DocumentSource struct {
	Location string `json:"location"` // local(本地磁盘), cloud(云, 如 oss url), db(就存储在 db)
	Path     string `json:"path"`
	Body     string `json:"body"` // set if location = db
}

type Document struct {
	Id     int64 `json:"id"`
	BookId int64 `json:"book_id"`
	//CategoryId      int64           `json:"category_id"`
	Title           string          `json:"title"`
	Source          DocumentSource  `json:"source" xorm:"json"`
	EmbeddingStatus EmbeddingStatus `json:"embedding_status"`
	CreatedAt       time.Time       `json:"created_at" xorm:"created"`
	UpdatedAt       time.Time       `json:"updated_at" xorm:"updated"`

	Fragments []Fragment `json:"fragments" xorm:"-"`
}

// Fragment 学习好的片段
type Fragment struct {
	Id         int64     `json:"id"`
	DocumentId int64     `json:"document_id"`
	BookId     int64     `json:"book_id"`
	Body       string    `json:"body"`
	StartIndex int       `json:"start_index"`
	EndIndex   int       `json:"end_index"`
	Vector     []float32 `json:"vector"`
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"updated"`
	Md5        string    `json:"md5"`
	Similarity float64   `json:"similarity" xorm:"->"`
}

type TaskLog struct {
	Id         string    `json:"id"`
	TargetId   string    `json:"target_id"`
	TargetType string    `json:"target_type"` // article
	Process    string    `json:"process"`     // e.g. 1/2
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type EmbeddingStatus = string

const (
	Pendding   = "pendding"
	Processing = "processing"
	Success    = "success"
	Retry      = "retry"
	Fail       = "fail"
)
