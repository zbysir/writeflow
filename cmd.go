package writeflow

import (
	"context"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/pkg/schema"
	"io/fs"
)

type funCMD struct {
	f interface{}
}

func (f *funCMD) Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	return execFunc(ctx, f.f, params)
}
func (f *funCMD) Schema(ctx context.Context) schema.CMDSchema {
	return schema.CMDSchema{}
}

func FunCMD(fun interface{}) schema.CMDer {
	return &funCMD{f: fun}
}

type GoPkgCMD struct {
	fs           fs.FS
	goPath       string // ./_pkg
	packagePath  string // examplegocmd
	execFuncName string // examplegocmd.Exec

	innerCMD schema.CMDer
}

type NewCmd func(config map[string]interface{}) (schema.CMDer, error)

func NewGoPkgCMD(fs fs.FS, goPath string, packagePath string, execFuncName string) (*GoPkgCMD, error) {
	i := interp.New(interp.Options{
		GoPath:               goPath,
		SourcecodeFilesystem: fs,
	})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, err
	}

	_, err = i.Eval(fmt.Sprintf(`import "%s"`, packagePath))
	if err != nil {
		return nil, err
	}

	res, err := i.Eval(execFuncName)
	if err != nil {
		return nil, err
	}

	fn := res.Interface().(NewCmd)

	inner, err := fn(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &GoPkgCMD{innerCMD: inner}, nil
}

func (g *GoPkgCMD) Schema(ctx context.Context) schema.CMDSchema {
	return g.innerCMD.Schema(ctx)
}

func (g *GoPkgCMD) Exec(ctx context.Context, params []interface{}) (rsp []interface{}, err error) {
	return g.innerCMD.Exec(ctx, params)
}
