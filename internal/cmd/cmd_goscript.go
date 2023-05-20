package cmd

import (
	"context"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/pkg/schema"
	"io/fs"
)

type GoScriptCMD struct {
	innerCMD schema.CMDer
}

// src:
// package examplegocmd
// import "github.com/zbysir/writeflow/pkg/schema"
// func Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {}
// func Schema() schema.Schema {}

func NewGoScript(fs fs.FS, goPath string, src string) (*GoScriptCMD, error) {
	i := interp.New(interp.Options{
		GoPath:               goPath,
		SourcecodeFilesystem: fs,
	})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, err
	}

	err = i.Use(Symbols())
	if err != nil {
		return nil, err
	}

	_, err = i.Eval(src)
	if err != nil {
		return nil, fmt.Errorf("failed to eval import: %w", err)
	}

	execFun, err := i.Eval("main.Exec")
	if err != nil {
		return nil, fmt.Errorf("failed to eval Exec: %w", err)
	}

	config := execFun.Interface()

	inner := config.(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error))
	return &GoScriptCMD{innerCMD: ExecFun(inner)}, nil
}

func (g *GoScriptCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return g.innerCMD.Exec(ctx, params)
}
