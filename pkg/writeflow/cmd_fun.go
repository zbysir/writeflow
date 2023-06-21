package writeflow

import (
	"context"
)

type ExecFun func(ctx context.Context, params Map) (rsp Map, err error)

func (e ExecFun) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	return e(ctx, params)
}

func NewFun(fun func(ctx context.Context, params Map) (rsp Map, err error)) CMDer {
	return ExecFun(fun)
}
