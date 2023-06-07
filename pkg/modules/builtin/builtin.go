package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/modules"
	"github.com/zbysir/writeflow/pkg/schema"
	"reflect"
	"regexp"
	"strings"
)

type Builtin struct {
}

func (b *Builtin) GoSymbols() map[string]map[string]reflect.Value {
	return nil
}

func New() *Builtin {
	return &Builtin{}
}

var _ modules.Module = (*Builtin)(nil)

func (b *Builtin) Info() modules.ModuleInfo {
	return modules.ModuleInfo{
		NameSpace: "builtin",
	}
}

// 以下组件 key 需要前端特殊处理：
// go_script: 有个编辑器直接编辑 go 代码
// output: 可以按照格式（如 markdown）显示输出

func (b *Builtin) Categories() []model.Category {
	return []model.Category{
		{
			Key: "input",
			Name: map[string]string{
				"zh-CN": "输入",
			},
			Desc: nil,
		},
		{
			Key: "output",
			Name: map[string]string{
				"zh-CN": "输出",
			},
			Desc: nil,
		},
		{
			Key: "data",
			Name: map[string]string{
				"zh-CN": "数据",
			},
			Desc: nil,
		},
		{
			Key: "text",
			Name: map[string]string{
				"zh-CN": "文本处理",
			},
			Desc: nil,
		},
		{
			Key: "script",
			Name: map[string]string{
				"zh-CN": "脚本",
			},
			Desc: nil,
		},
	}
}

func (b *Builtin) Components() []model.Component {
	return []model.Component{
		{
			Id:       0,
			Key:      "input_string",
			Category: "input",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "输入字符串",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "raw",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputParams: []model.NodeInputParam{
					{
						Id: "",
						Name: map[string]string{
							"zh-CN": "字符串",
						},
						Key:      "default",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []model.NodeAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:      "default",
						Type:     "string",
						List:     false,
						Optional: false,
					},
				},
			},
		},
		{
			Id:       0,
			Key:      "go_script",
			Category: "script",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "Golang 脚本",
				},
				Source: model.ComponentSource{
					CmdType:    model.GoScriptCmd,
					BuiltinCmd: "",
					GoPackage:  model.ComponentGoPackage{},
					GoScript: model.ComponentGoScript{
						Script: `package main
import (
    "context"
    "fmt"
)

func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
    return params, nil
}
					`,
					},
				},
				InputParams: []model.NodeInputParam{
					{
						Id: "",
						Name: map[string]string{
							"zh-CN": "字符串",
						},
						Key:      "default",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []model.NodeAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:      "default",
						Type:     "string",
						List:     false,
						Optional: false,
					},
				},
			},
		},
		{
			Id:       0,
			Key:      "output",
			Category: "output",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "输出",
				},
				Source: model.ComponentSource{
					CmdType:   model.NothingCmd,
					GoPackage: model.ComponentGoPackage{},
					GoScript:  model.ComponentGoScript{},
				},
				InputParams: []model.NodeInputParam{
					{
						Id: "",
						Name: map[string]string{
							"zh-CN": "数据",
						},
						Key:      "default",
						Type:     "any",
						Optional: true,
					},
				},
			},
		},
		{
			Id:       0,
			Key:      "select",
			Category: "data",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "通过路径选择数据",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "select",
				},
				InputAnchors: []model.NodeAnchor{
					{
						Name: map[string]string{
							"zh-CN": "数据",
						},
						Key:      "data",
						Type:     "any",
						Optional: true,
					},
				},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "Path",
						},
						Key:      "path",
						Type:     "string",
						Optional: true,
					},
				},
			},
		},
		{
			Id:       0,
			Key:      "record",
			Category: "data",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "组合多个数据为一个集合",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "record",
				},
				DynamicInput: true,
				// dynamic input
				InputAnchors: []model.NodeAnchor{},
				InputParams:  []model.NodeInputParam{},
				OutputAnchors: []model.NodeAnchor{
					{
						Name: map[string]string{
							"zh-CN": "集合",
						},
						Key:      "default",
						Type:     "any",
						Optional: false,
					},
				},
			},
		},
		{
			Key:      "template_text",
			Category: "text",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "模板文本",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "template_text",
				},

				// dynamic input
				DynamicInput: true,
				InputAnchors: []model.NodeAnchor{},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "模板",
						},
						Key:  "template",
						Type: "string",
					},
				},
				OutputAnchors: []model.NodeAnchor{
					{
						Key:  "default",
						Type: "string",
					},
				},
			},
		},
	}
}

func (b *Builtin) Cmd() map[string]schema.CMDer {
	return map[string]schema.CMDer{
		// 原封不动的返回节点入参
		"raw": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			//log.Infof("raw params: %+v", params)
			return params, nil
		}),
		"record": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			return map[string]interface{}{"default": params}, nil
		}),
		"template_text": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			tpl := params["template"].(string)

			vars, _ := json.Marshal(params)
			var varMap map[string]interface{}
			err = json.Unmarshal(vars, &varMap)
			if err != nil {
				return nil, err
			}
			// match {{abc}}
			reg := regexp.MustCompile(`{{\s*(\w+)\s*}}`)
			s := reg.ReplaceAllStringFunc(tpl, func(s string) string {
				s = strings.Trim(s, "{} ")
				r, ok := lookupMap(varMap, strings.Split(s, ".")...)
				if !ok {
					return s
				}
				return fmt.Sprintf("%v", r)
			})

			return map[string]interface{}{"default": s}, nil
		}),
		// 通过路径选择入参返回
		"select": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			//log.Infof("select params: %+v", params)
			p := params["path"]
			if p == nil {
				return map[string]interface{}{"default": params}, nil
			}
			d := params["data"]
			if d == nil {
				return map[string]interface{}{"default": nil}, nil
			}

			path := p.(string)
			data := d.(map[string]interface{})
			// TODO 使用 goja 来实现，goja 支持更多类型，并且语法和 js 一致。
			i, ok := lookupMap(data, strings.Split(path, ".")...)
			if !ok {
				return nil, nil
			}
			return map[string]interface{}{"default": i}, nil
		}),
	}
}

func lookupMap(m map[string]interface{}, keys ...string) (interface{}, bool) {
	var c interface{} = m
	for _, k := range keys {
		mm, ok := c.(map[string]interface{})
		if !ok {
			mmi, ok := c.(map[interface{}]interface{})
			if !ok {
				return nil, false
			}

			mm = make(map[string]interface{})
			for k, v := range mmi {
				mm[fmt.Sprintf("%v", k)] = v
			}
		}

		i, ok := mm[k]
		if !ok {
			return nil, false
		}

		c = i
	}

	return c, true
}
