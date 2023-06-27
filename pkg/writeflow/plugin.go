package writeflow

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/pkg/plugin"
	"github.com/zbysir/writeflow/pkg/writeflow/gosymbols"
	"io/fs"
	"reflect"
	"strings"
)

// GoPkgPluginManager list all plugin in a dir
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
	pkgName string // default is main
}

func NewGoPkgPlugin(fs fs.FS) *GoPkgPlugin {
	return &GoPkgPlugin{fs: fs}
}

type removePrefixFs struct {
	fs     fs.FS
	prefix string
}

func (p *removePrefixFs) Open(name string) (fs.File, error) {
	return p.fs.Open(strings.TrimPrefix(name, p.prefix))
}

func RemovePrefixFs(fs fs.FS, prefix string) fs.FS {
	return &removePrefixFs{
		fs:     fs,
		prefix: prefix,
	}
}

func (p *GoPkgPlugin) Register(r plugin.Register) (err error) {
	if p.pkgName == "" {
		p.pkgName = "main"
	}
	i := interp.New(interp.Options{
		GoPath: "./",
		// wrap src/pkgname
		SourcecodeFilesystem: RemovePrefixFs(p.fs, "src/"+p.pkgName),
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
