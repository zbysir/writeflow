package cmd

import (
	"context"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/pkg/schema"
	"io/fs"
	"reflect"
)

type GoPkgCMD struct {
	innerCMD schema.CMDer
}

type NewCmd func(config map[string]interface{}) (schema.CMDer, error)

type _CMDer struct {
	IValue  interface{}
	WExec   func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	WSchema func() schema.CMDSchema
}

func (p _CMDer) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return p.WExec(ctx, params)
}

func (p _CMDer) Schema() schema.CMDSchema {
	return p.WSchema()
}

func Symbols() map[string]map[string]reflect.Value {
	return map[string]map[string]reflect.Value{
		"github.com/zbysir/writeflow/pkg/schema/schema": {
			"CMDer":           reflect.ValueOf((*schema.CMDer)(nil)),
			"_CMDer":          reflect.ValueOf((*_CMDer)(nil)),
			"CMDSchema":       reflect.ValueOf((*schema.CMDSchema)(nil)),
			"CMDSchemaParams": reflect.ValueOf((*schema.CMDSchemaParams)(nil)),
		},
	}
}

func NewGoPkg(fs fs.FS, goPath string, packagePath string) (*GoPkgCMD, error) {
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

	_, err = i.Eval(fmt.Sprintf(`import "%v"`, packagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to eval import: %w", err)
	}

	configi, err := i.Eval(fmt.Sprintf("%v.CreateConfig()", packagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to eval CreateConfig: %w", err)
	}

	newFun, err := i.Eval(fmt.Sprintf("%v.New", packagePath))
	if err != nil {
		return nil, err
	}

	config := configi.Interface()

	inner := newFun.Call([]reflect.Value{reflect.ValueOf(config)})[0].Interface().(schema.CMDer)
	return &GoPkgCMD{innerCMD: inner}, nil
}

func (g *GoPkgCMD) Schema() schema.CMDSchema {
	return g.innerCMD.Schema()
}

func (g *GoPkgCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return g.innerCMD.Exec(ctx, params)
}
