package cmd

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/zbysir/gojsx"
)

type JavaScriptCMD struct {
	src string
}

// src:
// function exec(params){return params}

func NewJavaScript(src string) (*JavaScriptCMD, error) {
	return &JavaScriptCMD{src: src}, nil
}

func newJsRuntime(src string) (*goja.Runtime, goja.Callable, error) {
	r := goja.New()
	f, err := r.RunScript("java_script_cmd", fmt.Sprintf("(%s)", src))
	if err != nil {
		return nil, nil, fmt.Errorf("run javascript error: %w", err)
	}
	c, ok := gojsx.AssertFunction(f)
	if !ok {
		return nil, nil, fmt.Errorf("can't export 'exec' function in javascript")
	}

	return r, c, nil
}

func (g *JavaScriptCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	r, call, err := newJsRuntime(g.src)
	if err != nil {
		return nil, err
	}

	rspj, err := call(nil, r.ToValue(params))
	if err != nil {
		return nil, err
	}
	rsp = map[string]interface{}{}
	err = r.ExportTo(rspj, &rsp)
	if err != nil {
		return nil, fmt.Errorf("export javascript return to map[string]interface{}")
	}

	return
}
