package writeflow

import (
	plugin2 "github.com/zbysir/writeflow/pkg/plugin"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type sysFs struct {
	root string
}

func (s sysFs) Open(name string) (fs.File, error) {
	//os.Chdir()
	return os.Open(filepath.Join(s.root, name))
}

type mockRegister struct {
	ms []plugin2.Module
}

func (m2 *mockRegister) RegisterModule(m plugin2.Module) {
	m2.ms = append(m2.ms, m)
}

func TestGoPkgPlugin(t *testing.T) {
	p := GoPkgPlugin{
		fs:      sysFs{"/Users/bysir/goproj/bysir/writeflow-plugin-llm"},
		pkgName: "writeflow_plugin_llm",
	}
	r := &mockRegister{}
	err := p.Register(r)
	if err != nil {
		t.Fatal(err)
	}

	r.ms[0].Cmd()
	t.Logf("%+v", r.ms)
}
