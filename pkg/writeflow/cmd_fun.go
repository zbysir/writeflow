package writeflow

import (
	"context"
)

type ExecFun func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)

func (e ExecFun) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return e(ctx, params)
}

func NewFun(fun func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)) CMDer {
	return ExecFun(fun)
}
