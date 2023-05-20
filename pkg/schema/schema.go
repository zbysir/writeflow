package schema

import "context"

type CMDer interface {
	Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	//Schema() CMDSchema
}
