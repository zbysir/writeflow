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
	"io/fs"
	"os"
	"reflect"
	"strings"
)

// GoPkgPluginManager list all plugin in a dir
type GoPkgPluginManager struct {
	cacheFs    billy.Filesystem // default is os.TempDir()
	pluginUrls []PluginSource   // e.g. ["https://github.com/zbysir/writeflow-plugin-llm"]
}

type PluginSource struct {
	Url    string
	Enable bool
}

func NewGoPkgPluginManager(cacheFs billy.Filesystem, pluginUrls []PluginSource) *GoPkgPluginManager {
	if cacheFs == nil {
		cacheFs = osfs.New(os.TempDir())
	}
	return &GoPkgPluginManager{cacheFs: cacheFs, pluginUrls: pluginUrls}
}

func (m *GoPkgPluginManager) Load() ([]*GoPkgPlugin, error) {
	ps := make([]*GoPkgPlugin, 0)

	// TODO 并行使用 git 下载插件
	for _, url := range m.pluginUrls {
		if !url.Enable {
			ps = append(ps, &GoPkgPlugin{
				Source:  url.Url,
				Enable:  false,
				fs:      nil,
				pkgName: "",
			})
			return nil, nil
		}

		dir := strings.TrimPrefix(url.Url, "https://")
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
		err = g.Pull(url.Url, "master", true)
		if err != nil {
			return nil, err
		}

		ps = append(ps, NewGoPkgPlugin(gobilly.NewStdFs(pluginFs), url.Url))
	}

	return ps, nil
}

type GoPkgPlugin struct {
	Source  string
	Enable  bool
	fs      fs.FS
	pkgName string // default is main
}

func NewGoPkgPlugin(fs fs.FS, source string) *GoPkgPlugin {
	return &GoPkgPlugin{fs: fs, Source: source, Enable: true}
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
