package modules

import (
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
)

type ModuleInfo struct {
	NameSpace string
}

type Module interface {
	Info() ModuleInfo // 模块名
	Categories() []model.Category
	Components() []model.Component
	Cmd() map[string]schema.CMDer
}
