package writeflow

import (
	"context"
)

type funCmd struct {
	f func(ctx context.Context, params Map) (rsp Map, err error)
}

func (f *funCmd) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	return f.f(ctx, params)
}

func NewFun(fun func(ctx context.Context, params Map) (rsp Map, err error)) CMDer {
	return ExecFun(fun)
}

func NewFunMap(fun func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)) CMDer {
	return ExecFunMap(fun)
}
