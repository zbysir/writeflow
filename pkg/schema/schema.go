package schema

import "context"

type CMDSchemaParams struct {
	Key         string
	Type        string
	NameLocales string
	DescLocales string
}

type CMDSchema struct {
	Inputs      []CMDSchemaParams
	Outputs     []CMDSchemaParams
	Key         string
	NameLocales string
	DescLocales string
}

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	Schema() CMDSchema
}
