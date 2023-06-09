package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/pkg/modules"
	"github.com/zbysir/writeflow/pkg/schema"
	"github.com/zbysir/writeflow/pkg/writeflow"
	"io/ioutil"
	"net/http"
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
			Key: "logic",
			Name: map[string]string{
				"zh-CN": "逻辑",
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
			Type:     "input_string",
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
						Name: map[string]string{
							"zh-CN": "字符串",
						},
						Key:      "default",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "string",
						List: false,
					},
				},
			},
		},
		{
			Id:       0,
			Type:     "call_http",
			Category: "data",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "HTTP 请求",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "call_http",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "URL",
						},
						Key:      "url",
						Type:     "string",
						Optional: true,
					},
					{
						Name: map[string]string{
							"zh-CN": "Body",
						},
						Key:      "body",
						Type:     "any",
						Optional: true,
					},
					{
						Name: map[string]string{
							"zh-CN": "方法 [GET/POST/PUT/DELETE]",
						},
						Key:      "method",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "any",
						List: false,
					},
				},
			},
		},
		{
			Id:       0,
			Type:     "switch",
			Category: "logic",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "Switch",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "_switch",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				// DynamicInput for conditions
				DynamicInput: true,
				InputAnchors: []model.NodeInputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Data",
						},
						Key:  "data",
						Type: "any",
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "any",
					},
				},
				// switch 的 output 需要前端做特殊处理，output keys 需要等于 input keys
			},
		},
		{
			Id:       0,
			Type:     "for",
			Category: "logic",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "For",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "_for",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputAnchors: []model.NodeInputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Data",
						},
						Key:  "data",
						Type: "any",
					},
					{
						Name: map[string]string{
							"zh-CN": "Item",
						},
						Key:  "item",
						Type: "any",
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "any",
					},
				},
			},
		},
		{
			Id:       0,
			Type:     "go_script",
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
						InputKey: "script",
					},
				},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "脚本代码",
						},
						Key:         "script",
						Type:        "string",
						DisplayType: "code/go",
						Optional:    false,
					},
				},
				Inputs: map[string]string{
					"script": `package main
import (
    "context"
    "fmt"
)

func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
    return params, nil
}
`,
				},
				DynamicInput:  true,
				OutputAnchors: []model.NodeOutputAnchor{},
			},
		},
		{
			Id:       0,
			Type:     "output",
			Category: "output",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "输出",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "raw",
				},
				InputAnchors: []model.NodeInputAnchor{
					{
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
			Type:     "select",
			Category: "data",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "通过路径选择数据",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "select",
				},
				InputAnchors: []model.NodeInputAnchor{
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
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "any",
					},
				},
			},
		},
		{
			Id:       0,
			Type:     "record",
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
				InputAnchors: []model.NodeInputAnchor{},
				InputParams:  []model.NodeInputParam{},
				OutputAnchors: []model.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "集合",
						},
						Key:  "default",
						Type: "any",
					},
				},
			},
		},
		{
			Type:     "template_text",
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
				InputAnchors: []model.NodeInputAnchor{},
				InputParams: []model.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "模板",
						},
						Key:  "template",
						Type: "string",
					},
				},
				OutputAnchors: []model.NodeOutputAnchor{
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
			l, err := writeflow.LookInterface(d, path)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"default": l}, nil
		}),
		// 多个分支，只会有一个 key 为 true，如果所有的 key 都是 false，则 default 为 true
		"switch": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			d := params["data"]
			if d == nil {
				return map[string]interface{}{"default": true}, nil
			}

			rsp = map[string]interface{}{}

			// 根据 input 定义顺序执行
			keys := writeflow.GetInputKeys(ctx)
			if len(keys) == 0 {
				for k := range params {
					if k == "data" {
						continue
					}
					keys = append(keys, k)
				}
			}

			// 除了 data 都是条件
			for _, k := range keys {
				if k == "data" {
					continue
				}

				v := params[k]
				// data.a
				condition, ok := v.(string)
				if !ok {
					log.Errorf("switch params error: %+v is not string", params)
					continue
				}

				l, err := writeflow.LookInterface(d, condition)
				if err != nil {
					return nil, err
				}

				if l != nil && l != false {
					rsp[k] = true
					break
				}
			}

			if len(rsp) == 0 {
				// 如果没有进入任何分支，则会进入 default 分支
				rsp["default"] = true
			}

			return rsp, nil
		}),
		"call_http": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			p := params["url"]
			if p == nil {
				return map[string]interface{}{"default": params}, nil
			}
			d := params["body"]
			if d == nil {
				return map[string]interface{}{"default": nil}, nil
			}

			m := params["method"]
			if d == nil {
				return map[string]interface{}{"default": nil}, nil
			}

			path := p.(string)
			data := d.(string)
			method := m.(string)
			httpClient := &http.Client{}
			req, err := http.NewRequest(method, path, bytes.NewBuffer([]byte(data)))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := httpClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var i interface{}
			err = json.Unmarshal(body, &i)
			if err != nil {
				return nil, err
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
