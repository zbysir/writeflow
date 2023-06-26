package writeflow

import (
	"context"
)

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
}

type SchemaParams struct {
	Key  string `json:"key" yaml:"key"`
	Type string `json:"type" yaml:"type"` // literal, anchor
	//Literal string `json:"literal"` // 字面量
	NameLocales map[string]string `json:"name_locales,omitempty" yaml:"name_locales"`
	DescLocales map[string]string `json:"desc_locales,omitempty" yaml:"desc_locales"`
}

type Schema struct {
	Key         string            `json:"key" yaml:"key"`
	Inputs      []SchemaParams    `json:"inputs" yaml:"inputs"`
	Outputs     []SchemaParams    `json:"outputs" yaml:"outputs"`
	NameLocales map[string]string `json:"name_locales,omitempty" yaml:"name_locales"`
	DescLocales map[string]string `json:"desc_locales,omitempty" yaml:"desc_locales"`
}
