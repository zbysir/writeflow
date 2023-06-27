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
	return os.Open(filepath.Join(s.root, name))
}

type mockRegister struct {
	ms []plugin2.Plugin
}

func (m2 *mockRegister) RegisterPlugin(m plugin2.Plugin) {
	m2.ms = append(m2.ms, m)
}

func TestGoPkgPlugin(t *testing.T) {
	p := NewGoPkgPlugin(NewSysFs("/Users/bysir/goproj/bysir/writeflow-plugin-llm"))
	r := &mockRegister{}
	err := p.Register(r)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", r.ms[0].Cmd())
}
