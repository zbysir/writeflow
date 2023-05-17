package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zbysir/writeflow/pkg/schema"
	"reflect"
	"strconv"
)

type funCMD struct {
	f      interface{}
	schema schema.CMDSchema
}

func (f *funCMD) Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return execFunc(ctx, f.f, params)
}
func (f *funCMD) Schema() schema.CMDSchema {
	return f.schema
}

func (f *funCMD) SetSchema(s schema.CMDSchema) *funCMD {
	f.schema = s
	return f
}

func NewFun(fun interface{}) *funCMD {
	return &funCMD{f: fun}
}

func execFunc(ctx context.Context, fun interface{}, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	if xfun, ok := fun.(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error)); ok {
		return xfun(ctx, params)
	}

	callParams := []reflect.Value{reflect.ValueOf(ctx)}

	for _, p := range params {
		callParams = append(callParams, reflect.ValueOf(p))
	}
	funv := reflect.ValueOf(fun)

	ty := funv.Type().NumIn()
	for i := 0; i < ty; i++ {
		wantp := funv.Type().In(i)
		inp := callParams[i].Type()

		//fmt.Printf("wantp:%v, inp:%v %v\n", wantp.String(), inp.String(), inp.AssignableTo(wantp))

		// TODO 如果目标是数组，则使用 Append 而不是直接赋值，来源可以支持数组 Item

		// 如果类型不匹配，则尝试通过 json 转换
		if !inp.AssignableTo(wantp) {
			bs, _ := json.Marshal(callParams[i].Interface())
			w := reflect.New(wantp)
			err = json.Unmarshal(bs, w.Interface())
			if err != nil {
				return nil, fmt.Errorf("can not convert %v to %v, err: %w", inp.String(), wantp.String(), err)
			}

			callParams[i] = w.Elem()
		}

		//if wantp.String() == "[]string" && inp.String() == "[]interface {}" {
		//	callParams[i] = reflect.ValueOf(interfaceTo[string](callParams[i].Interface().([]interface{})))
		//}

	}

	rv := funv.Call(callParams)
	rsp = map[string]interface{}{}
	var rerr error
	l := len(rv)
	for i, v := range rv {
		if i == l-1 {
			last := v
			switch last.Kind() {
			case reflect.Interface:
				err, ok := last.Interface().(error)
				if ok {
					rerr = err
					continue
				}
			}
			rsp[strconv.Itoa(len(rsp))] = v.Interface()
		} else {
			rsp[strconv.Itoa(len(rsp))] = v.Interface()
		}
	}

	return rsp, rerr
}
