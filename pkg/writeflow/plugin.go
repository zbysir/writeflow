package writeflow

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	plugin2 "github.com/zbysir/writeflow/pkg/plugin"
	"github.com/zbysir/writeflow/pkg/writeflow/gosymbols"
	"io/fs"
	"reflect"
	"strings"
)

type GoPkgPluginManager struct {
	fs fs.FS
}

func (m *GoPkgPluginManager) Load() ([]GoPkgPlugin, error) {
	ds, err := fs.ReadDir(m.fs, ".")
	if err != nil {
		return nil, err
	}

	ps := make([]GoPkgPlugin, 0)
	for _, d := range ds {
		if d.IsDir() {
			sub, err := fs.Sub(m.fs, d.Name())
			if err != nil {
				return nil, err
			}
			ps = append(ps, GoPkgPlugin{
				fs:      sub,
				pkgName: d.Name(),
			})
		}
	}

	return ps, nil
}

type GoPkgPlugin struct {
	fs      fs.FS
	pkgName string
}

// type GoPkgPluginModule struct{
//
// }

type AddPrefixFs struct {
	fs     fs.FS
	prefix string
}

func (p *AddPrefixFs) Open(name string) (fs.File, error) {
	return p.fs.Open(strings.TrimPrefix(name, p.prefix))
}

func NewRemovePrefixFs(fs fs.FS, prefix string) *AddPrefixFs {
	return &AddPrefixFs{
		fs:     fs,
		prefix: prefix,
	}
}

func (p *GoPkgPlugin) Register(r plugin2.ModuleRegister) (err error) {
	i := interp.New(interp.Options{
		GoPath: "./",
		// wrap src/pkgname
		SourcecodeFilesystem: NewRemovePrefixFs(p.fs, "src/"+p.pkgName),
	})

	err = i.Use(stdlib.Symbols)
	if err != nil {
		return err
	}

	err = i.Use(gosymbols.Symbols)
	if err != nil {
		return err
	}

	_, err = i.Eval(fmt.Sprintf(`import "%v"`, p.pkgName))
	if err != nil {
		return fmt.Errorf("failed to eval import: %w", err)
	}

	// func Register(r ModuleRegister){}
	newFun, err := i.Eval(fmt.Sprintf("%v.Register", p.pkgName))
	if err != nil {
		return err
	}

	_ = newFun.Call([]reflect.Value{reflect.ValueOf(r)})

	return nil
}
