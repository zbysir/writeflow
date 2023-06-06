package modules

import (
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
	"reflect"
)

type ModuleInfo struct {
	NameSpace string
}

type Module interface {
	Info() ModuleInfo
	Categories() []model.Category
	Components() []model.Component
	Cmd() map[string]schema.CMDer
	GoSymbols() map[string]map[string]reflect.Value // If you want to create a component by go code, you need to implement this method. more detail see `yaegi` project.
}
