package builtin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/pkg/schema"
	"github.com/zbysir/writeflow/pkg/writeflow"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Builtin struct {
}

func (b *Builtin) GoSymbols() map[string]map[string]reflect.Value {
	return nil
}

func New() *Builtin {
	return &Builtin{}
}

var _ writeflow.Module = (*Builtin)(nil)

func (b *Builtin) Info() writeflow.ModuleInfo {
	return writeflow.ModuleInfo{
		NameSpace: "builtin",
	}
}

// 以下组件 key 需要前端特殊处理：
// go_script: 有个编辑器直接编辑 go 代码
// output: 可以按照格式（如 markdown）显示输出

func (b *Builtin) Categories() []writeflow.Category {
	return []writeflow.Category{
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

func (b *Builtin) Components() []writeflow.Component {
	return []writeflow.Component{
		{
			Id:       0,
			Type:     "input_string",
			Category: "input",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "输入字符串",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "raw",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "字符串",
						},
						Key:      "default",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Type:     "input_string_password",
			Category: "input",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Password",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "raw",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "字符串",
						},
						Key:         "default",
						Type:        "string",
						DisplayType: "password",
						Optional:    true,
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Type:     "params",
			Category: "input",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "请求参数",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "_params",
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Type:     "call_http",
			Category: "data",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "HTTP 请求",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "call_http",
				},
				InputParams: []writeflow.NodeInputParam{
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
							"zh-CN": "方法 [GET/POST/PUT/DELETE]",
						},
						Key:      "method",
						Type:     "string",
						Optional: true,
					},
					{
						InputType: writeflow.NodeInputAnchor,

						Name: map[string]string{
							"zh-CN": "Body",
						},
						Key:      "body",
						Type:     "any",
						Optional: true,
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Switch",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "_switch",
				},
				// DynamicInput for conditions
				DynamicInput: true,
				InputParams: []writeflow.NodeInputParam{
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "Data",
						},
						Key:  "data",
						Type: "any",
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Type:     "for",
			Category: "logic",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "For",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "_for",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "Data",
						},
						Key:  "data",
						Type: "any",
					},
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "Item",
						},
						Key:  "item",
						Type: "any",
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
					{
						Name: map[string]string{
							"zh-CN": "Item",
						},
						Key:  "item",
						Type: "any",
					},
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
			Type:     "sleep",
			Category: "logic",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Sleep",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "sleep",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "Second",
						},
						Key:  "second",
						Type: "number",
					},
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "Default",
						},
						Key:  "default",
						Type: "any",
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Golang 脚本",
				},
				Source: writeflow.ComponentSource{
					CmdType: writeflow.GoScriptCmd,
				},
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "脚本代码",
						},
						Key:         "script",
						Type:        "string",
						DisplayType: "code/go",
						Optional:    false,
						Value: `package main
import (
    "context"
    "strings"
)

func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
	return map[string]interface{}{"default": strings.TrimSpace(params["default"].(string))} , nil
}
`,
					},
				},
				DynamicInput:  true,
				DynamicOutput: true,
				OutputAnchors: []writeflow.NodeOutputAnchor{},
			},
		},
		{
			Id:       0,
			Type:     "js_script",
			Category: "script",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Javascript 脚本",
				},
				Source: writeflow.ComponentSource{
					CmdType: writeflow.JavaScriptCmd,
				},
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "脚本代码",
						},
						Key:         "script",
						Type:        "string",
						DisplayType: "code/javascript",
						Optional:    false,
						Value:       `function exec (params){return params}`,
					},
				},
				DynamicInput:  true,
				DynamicOutput: true,
			},
		},
		{
			Id:       0,
			Type:     "output",
			Category: "output",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Output",
				},
				Description: map[string]string{
					"zh-CN": "输出",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "raw",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						InputType: writeflow.NodeInputLiteral,
						Name: map[string]string{
							"zh-CN": "Enable",
						},
						Key:      "_enable",
						Type:     "bool",
						Value:    true,
						Optional: true,
					},
					{
						InputType: writeflow.NodeInputAnchor,
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
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "Select",
				},
				Description: map[string]string{
					"zh-CN": "通过路径选择数据",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "select",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "数据",
						},
						Key:      "data",
						Type:     "any",
						Optional: true,
					},
					{
						Name: map[string]string{
							"zh-CN": "Path",
						},
						Key:      "path",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Type:     "list",
			Category: "data",
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "List",
				},
				Description: map[string]string{
					"zh-CN": "多个数据合成数组",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "list",
				},
				InputParams: []writeflow.NodeInputParam{
					{
						InputType: writeflow.NodeInputAnchor,
						Name: map[string]string{
							"zh-CN": "数据",
						},
						Key:      "data",
						Type:     "any",
						Optional: true,
						List:     true,
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "组合多个数据为一个集合",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "record",
				},
				DynamicInput: true,
				// dynamic input
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
			Data: writeflow.ComponentData{
				Name: map[string]string{
					"zh-CN": "模板文本",
				},
				Source: writeflow.ComponentSource{
					CmdType:    writeflow.BuiltInCmd,
					BuiltinCmd: "template_text",
				},

				// dynamic input
				DynamicInput: true,
				InputParams: []writeflow.NodeInputParam{
					{
						Name: map[string]string{
							"zh-CN": "模板",
						},
						Key:         "template",
						Type:        "string",
						DisplayType: "textarea",
					},
				},
				OutputAnchors: []writeflow.NodeOutputAnchor{
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
		"sleep": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			//log.Infof("raw params: %+v", params)
			s := cast.ToInt(params["second"])
			if s != 0 {
				time.Sleep(time.Duration(s) * time.Second)
			}
			return params, nil
		}),
		"record": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			return map[string]interface{}{"default": params}, nil
		}),
		"template_text": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			tpl := params["template"].(string)

			// match {{abc}}
			reg := regexp.MustCompile(`{{.+?}}`)
			s := reg.ReplaceAllStringFunc(tpl, func(s string) string {
				if err != nil {
					return ""
				}
				s = strings.TrimPrefix(s, "{{")
				s = strings.TrimSuffix(s, "}}")
				s = strings.TrimSpace(s)
				r, e := writeflow.LookInterface(params, s)
				if e != nil {
					err = fmt.Errorf("exec template exp '%s' error: %w", s, e)
					return ""
				}
				return fmt.Sprintf("%v", r)
			})
			if err != nil {
				return nil, err
			}

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
			l, err := writeflow.LookInterface(map[string]interface{}{"data": d}, path)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"default": l}, nil
		}),
		"list": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			//log.Infof("list params: %+v", params)
			p := params["data"]
			if p == nil {
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{"default": p}, nil
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
