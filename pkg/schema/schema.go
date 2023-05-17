package schema

import "context"

type CMDSchemaParams struct {
	Key         string
	Type        string
	NameLocales map[string]string
	DescLocales map[string]string
}

type CMDSchema struct {
	Inputs      []CMDSchemaParams
	Outputs     []CMDSchemaParams
	Key         string
	Name        string
	NameLocales map[string]string
	DescLocales map[string]string
}

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	Schema() CMDSchema
}
