package gosymbols

import "reflect"

//go:generate yaegi extract github.com/tmc/langchaingo/llms
//go:generate yaegi extract github.com/tmc/langchaingo/llms/openai
var Symbols = map[string]map[string]reflect.Value{}
