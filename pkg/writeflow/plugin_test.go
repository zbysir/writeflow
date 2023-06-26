package writeflow

import (
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
	ms []Module
}

func (m2 *mockRegister) RegisterModule(m Module) {
	m2.ms = append(m2.ms, m)
}

func TestGoPkgPlugin(t *testing.T) {
	p := GoPkgPlugin{
		fs:      sysFs{"/Users/bysir/goproj/bysir/writeflow-plugin-llm"},
		pkgName: "writeflow_plugin_llm",
	}
	r := &mockRegister{}
	err := p.Register(r)
	if ErrNodeUnreachable != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", r.ms)
}
