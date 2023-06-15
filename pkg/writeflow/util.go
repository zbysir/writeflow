package writeflow

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/zbysir/gojsx"
)

// LookInterface i: {"a": 1, "b": {"d": 2}}, support key: a, b.d
func LookInterface(i map[string]interface{}, key string) (v interface{}, err error) {
	r := goja.New()
	for k, v := range i {
		err = r.Set(k, v)
		if err != nil {
			return nil, err
		}
	}

	out, err := r.RunScript("look_interface", fmt.Sprintf("%v", key))
	if err != nil {
		return nil, gojsx.PrettifyException(err)
	}
	return out.Export(), nil
}

func ForInterface(i interface{}, n func(i interface{})) (err error) {
	r := goja.New()
	r.Set("data", i)
	r.Set("n", n)

	_, err = r.RunScript("for_interface", fmt.Sprintf("data.forEach(n)"))
	if err != nil {
		return gojsx.PrettifyException(err)
	}
	return nil
}
