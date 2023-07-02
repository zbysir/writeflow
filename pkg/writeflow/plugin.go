package writeflow

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/zbysir/writeflow/internal/pkg/git"
	"github.com/zbysir/writeflow/internal/pkg/gobilly"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/pkg/export"
	"github.com/zbysir/writeflow/pkg/writeflow/gosymbols"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// GoPkgPluginManager list all plugin in a dir
type GoPkgPluginManager struct {
	cacheFs billy.Filesystem // default is os.TempDir() + "/writeflow"
}

type PluginSource struct {
	Url string
}

func NewGoPkgPluginManager(cacheFs billy.Filesystem) *GoPkgPluginManager {
	if cacheFs == nil {
		//log.Infof("cachedir: %s", os.TempDir())
		cacheFs = osfs.New(filepath.Join(os.TempDir(), "writeflow"))
	}
	return &GoPkgPluginManager{cacheFs: cacheFs}
}

func (m *GoPkgPluginManager) Load(sourceUrl string) (*GoPkgPlugin, error) {
	dir := strings.TrimPrefix(sourceUrl, "https://")
	dir = strings.TrimSuffix(dir, ".git")
	err := m.cacheFs.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	// git clone
	pluginFs := chroot.New(m.cacheFs, dir)
	g, err := git.NewGit("", pluginFs, log.New(log.Options{
		IsDev:         false,
		To:            nil,
		DisableTime:   false,
		DisableLevel:  false,
		DisableCaller: true,
		CallerSkip:    0,
		Name:          "[Plugin]",
	}))
	if err != nil {
		return nil, err
	}

	// TODO Support specify branch
	err = g.Pull(sourceUrl, "master", true)
	if err != nil {
		return nil, err
	}
	return NewGoPkgPlugin(gobilly.NewStdFs(pluginFs), sourceUrl), nil
}

type GoPkgPlugin struct {
	Source  string
	fs      fs.FS
	pkgName string // default is main
}

func NewGoPkgPlugin(fs fs.FS, source string) *GoPkgPlugin {
	return &GoPkgPlugin{fs: fs, Source: source}
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

func (p *GoPkgPlugin) Register(r export.Register) (err error) {
	if p.pkgName == "" {
		p.pkgName = "main"
	}

	modFile, err := p.fs.Open("go.mod")
	if err != nil {
		return err
	}
	defer modFile.Close()
	modBytes, err := io.ReadAll(modFile)
	if err != nil {
		return err
	}
	pkgName := strings.Split(string(modBytes), "\n")[0]
	p.pkgName = strings.TrimPrefix(pkgName, "module ")

	//p.pkgName = "github.com/zbysir/writeflow_plugin_llm"

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

	// github.com/zbysir/writeflow_plugin_llm
	_, err = i.Eval(fmt.Sprintf(`import plugin "%v"`, p.pkgName))
	if err != nil {
		return fmt.Errorf("failed to eval import: %w", err)
	}

	// func Register(r ModuleRegister){}
	newFun, err := i.Eval(fmt.Sprintf("%v.Register", "plugin"))
	if err != nil {
		return fmt.Errorf("failed to eval Register: %w", err)
	}

	_ = newFun.Call([]reflect.Value{reflect.ValueOf(r)})

	return nil
}
