package explore

import (
	"context"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type CMDer interface {
	Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error)
}

type funCMD struct {
	f interface{}
}

func (f *funCMD) Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	return execFunc(ctx, f.f, params)
}

func FunCMD(fun interface{}) CMDer {
	return &funCMD{f: fun}
}

type GoCMD struct {
	script string
}

func (g GoCMD) Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	i := interp.New(interp.Options{
		GoPath: "./_pkg",
	})

	i.Use(stdlib.Symbols)
	_, err = i.Eval(`import "examplegocmd"`)
	if err != nil {
		panic(err)
	}

	res, err := i.Eval("examplegocmd.Exec")
	if err != nil {
		return nil, err
	}

	fn := res.Interface().(func(ctx context.Context, params []interface{}) (rsp []interface{}, err error))

	//log.Printf("res: %v", fn)

	return fn(ctx, params)
	//return nil, err
}
