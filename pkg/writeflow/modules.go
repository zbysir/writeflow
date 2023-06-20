package writeflow

import (
	"reflect"
)

type ModuleInfo struct {
	NameSpace string
}

type Module interface {
	Info() ModuleInfo
	Categories() []Category
	Components() []Component
	Cmd() map[string]CMDer
	GoSymbols() map[string]map[string]reflect.Value // If you want to create a component by go code, you need to implement this method. more detail see `yaegi` project.
}
