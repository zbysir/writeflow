package schema

import "context"

type CMDSchemaParams struct {
	Key  string
	Type string
	Desc string
}

type CMDSchema struct {
	Inputs  []CMDSchemaParams
	Outputs []CMDSchemaParams
	Name    string
	Desc    string
}

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	Schema(ctx context.Context) CMDSchema
}
