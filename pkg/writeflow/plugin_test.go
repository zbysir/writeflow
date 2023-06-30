package writeflow

import (
	"github.com/zbysir/writeflow/pkg/export"
	"testing"
)

type mockRegister struct {
	ms []export.Plugin
}

func (m2 *mockRegister) RegisterPlugin(m export.Plugin) {
	m2.ms = append(m2.ms, m)
}
func TestNewGoPkgPluginManager(t *testing.T) {
	pm := NewGoPkgPluginManager(nil)
	gg, err := pm.Load("https://github.com/zbysir/writeflow-plugin-llm")
	if err != nil {
		t.Fatal(err)
	}

	r := &mockRegister{}
	err = gg.Register(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoPkgPlugin(t *testing.T) {
	p := NewGoPkgPlugin(NewSysFs("/Users/bysir/goproj/bysir/writeflow-plugin-llm"), "")
	r := &mockRegister{}
	err := p.Register(r)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", r.ms[0].Cmd())
}
