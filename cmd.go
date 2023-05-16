package writeflow

import (
	"context"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/pkg/schema"
	"io/fs"
	"log"
	"reflect"
)

type funCMD struct {
	f interface{}
}

func (f *funCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return execFunc(ctx, f.f, params)
}
func (f *funCMD) Schema(ctx context.Context) schema.CMDSchema {
	return schema.CMDSchema{}
}

func FunCMD(fun interface{}) schema.CMDer {
	return &funCMD{f: fun}
}

type GoPkgCMD struct {
	innerCMD schema.CMDer
}

type NewCmd func(config map[string]interface{}) (schema.CMDer, error)

type _CMDer struct {
	IValue  interface{}
	WExec   func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)
	WSchema func(ctx context.Context) schema.CMDSchema
}

func (p _CMDer) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return p.WExec(ctx, params)
}

func (p _CMDer) Schema(ctx context.Context) schema.CMDSchema {
	return p.WSchema(ctx)
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

func NewGoPkgCMD(fs fs.FS, goPath string, packagePath string, execFuncName string) (*GoPkgCMD, error) {
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

	srcipt := fmt.Sprintf(`
package wrapper

import (
	"context"

	"%v"
	"github.com/zbysir/writeflow/pkg/schema"
)

func NewWrapper(c *%v.Config) (schema.CMDer, error) {
	p, err := %v.New(c)
	var pv schema.CMDer = p
	return pv, err
}
`, packagePath, packagePath, packagePath)

	log.Printf("script: %+v", srcipt)

	_, err = i.Eval(srcipt)
	if err != nil {
		return nil, fmt.Errorf("failed to eval wrapper: %w", err)
	}
	res, err := i.Eval("wrapper.NewWrapper")
	if err != nil {
		return nil, err
	}

	config := configi.Interface()

	inner := res.Call([]reflect.Value{reflect.ValueOf(config)})[0].Interface().(schema.CMDer)
	return &GoPkgCMD{innerCMD: inner}, nil
}

func (g *GoPkgCMD) Schema(ctx context.Context) schema.CMDSchema {
	return g.innerCMD.Schema(ctx)
}

func (g *GoPkgCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return g.innerCMD.Exec(ctx, params)
}
