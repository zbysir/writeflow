package schema

import "context"

type CMDSchemaParams struct {
	Key         string            `json:"key"`
	Type        string            `json:"type"`
	NameLocales map[string]string `json:"name_locales,omitempty"`
	DescLocales map[string]string `json:"desc_locales,omitempty"`
}

type CMDSchema struct {
	Inputs      []CMDSchemaParams `json:"inputs"`
	Outputs     []CMDSchemaParams `json:"outputs"`
	Key         string            `json:"key"`
	NameLocales map[string]string `json:"name_locales,omitempty"`
	DescLocales map[string]string `json:"desc_locales,omitempty"`
}

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	Schema() CMDSchema
}
