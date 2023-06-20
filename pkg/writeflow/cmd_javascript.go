package writeflow

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

func NewJavaScriptCMD(src string) (*JavaScriptCMD, error) {
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

func (g *JavaScriptCMD) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	r, call, err := newJsRuntime(g.src)
	if err != nil {
		return Map{}, err
	}

	rspj, err := call(nil, r.ToValue(params))
	if err != nil {
		return Map{}, err
	}
	rspr := map[string]interface{}{}
	err = r.ExportTo(rspj, &rsp)
	if err != nil {
		return Map{}, fmt.Errorf("export javascript return to map[string]interface{}")
	}

	rsp = NewMap(rspr)
	return
}
