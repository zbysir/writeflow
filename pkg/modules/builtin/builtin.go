package builtin

import (
	"context"
	"fmt"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/modules"
	"github.com/zbysir/writeflow/pkg/schema"
	"strings"
)

type Builtin struct {
}

func New() *Builtin {
	return &Builtin{}
}

var _ modules.Module = &Builtin{}

func (b *Builtin) Info() modules.ModuleInfo {
	return modules.ModuleInfo{
		NameSpace: "builtin",
	}
}

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
						Id:       "",
						Name:     nil,
						Key:      "default",
						Type:     "string",
						Optional: true,
					},
				},
				OutputAnchors: []model.NodeAnchor{
					{
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
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "raw",
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputParams: []model.NodeInputParam{
					{
						Id:       "",
						Name:     nil,
						Key:      "default",
						Type:     "string",
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
					GoPackage:  model.ComponentGoPackage{},
					GoScript:   model.ComponentGoScript{},
				},
				InputAnchors: []model.NodeAnchor{
					{
						Key:      "data",
						Type:     "any",
						List:     false,
						Optional: true,
					},
				},
				InputParams: []model.NodeInputParam{
					{
						Key:      "path",
						Type:     "string",
						Optional: true,
					},
				},
			},
		}, {
			Id:       0,
			Key:      "record",
			Category: "data",
			Data: model.ComponentData{
				Name: map[string]string{
					"zh-CN": "组合多个数据为一个集合",
				},
				Source: model.ComponentSource{
					CmdType:    model.BuiltInCmd,
					BuiltinCmd: "raw",
				},
				InputAnchors: []model.NodeAnchor{},
				InputParams:  []model.NodeInputParam{},
				OutputAnchors: []model.NodeAnchor{
					{
						Key:      "default",
						Type:     "any",
						List:     false,
						Optional: false,
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
