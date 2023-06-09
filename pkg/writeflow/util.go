package writeflow

import (
	"fmt"
	"github.com/dop251/goja"
)

func LookInterface(i interface{}, key string) (v interface{}, err error) {
	r := goja.New()
	r.Set("data", i)

	out, err := r.RunScript("look_interface", fmt.Sprintf("%v", key))
	if err != nil {
		return nil, err
	}
	return out.Export(), nil
}

func ForInterface(i interface{}, n func(i interface{})) (err error) {
	r := goja.New()
	r.Set("data", i)
	r.Set("n", n)

	_, err = r.RunScript("for_interface", fmt.Sprintf("data.forEach(n)"))
	if err != nil {
		return err
	}
	return nil
}
