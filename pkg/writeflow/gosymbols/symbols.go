package gosymbols

import "reflect"

//go:generate yaegi extract github.com/zbysir/writeflow/pkg/export
var Symbols = map[string]map[string]reflect.Value{}
