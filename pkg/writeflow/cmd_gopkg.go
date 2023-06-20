package writeflow

import (
	"context"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/fs"
	"reflect"
)

type GoPkgCMD struct {
	innerCMD CMDer
}

type NewCmd func(config map[string]interface{}) (CMDer, error)

type _CMDer struct {
	IValue interface{}
	WExec  func(ctx context.Context, params Map) (rsp Map, err error)
}

func (p _CMDer) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	return p.WExec(ctx, params)
}

func Symbols() map[string]map[string]reflect.Value {
	return map[string]map[string]reflect.Value{
		"github.com/zbysir/writeflow/pkg/schema/schema": {
			"CMDer":  reflect.ValueOf((*CMDer)(nil)),
			"_CMDer": reflect.ValueOf((*_CMDer)(nil)),
			//"Schema":       reflect.ValueOf((*schema.Schema)(nil)),
			//"SchemaParams": reflect.ValueOf((*schema.SchemaParams)(nil)),
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

	inner := newFun.Call([]reflect.Value{reflect.ValueOf(config)})[0].Interface().(CMDer)
	return &GoPkgCMD{innerCMD: inner}, nil
}

func (g *GoPkgCMD) Exec(ctx context.Context, params Map) (rsp Map, err error) {
	return g.innerCMD.Exec(ctx, params)
}
