package gosymbols

import "reflect"

//go:generate yaegi extract github.com/tmc/langchaingo/llms
//go:generate yaegi extract github.com/tmc/langchaingo/llms/openai
//go:generate yaegi extract github.com/zbysir/writeflow/pkg/plugin
var Symbols = map[string]map[string]reflect.Value{}
