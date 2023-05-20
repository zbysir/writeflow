package writeflow

import (
	"github.com/zbysir/writeflow/internal/model"
	"testing"
)

func TestComponentFromModel(t *testing.T) {
	m, err := ComponentFromModel(&model.Component{
		Id:  0,
		Key: "demo",
		Data: model.NodeData{
			Label:       "",
			Id:          "",
			Name:        "",
			Type:        "",
			Category:    "",
			Icon:        "",
			Description: "",
			Inputs:      nil,
			Source: model.NodeSource{
				Type:    "",
				CmdType: "go_script",
				GitUrl:  "",
				GoScript: `package main
					import (
				
				"context"
				"fmt"
				)
					
					func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
						return map[string]interface{}{"msg": fmt.Sprintf("hello %v, your age is: %v", params["name"], params["age"])}, nil
					}
					`,
			},
			InputAnchors: []model.NodeAnchor{
				{
					Id: "",
					Name: map[string]string{
						"zh-CN": "姓名",
					},
					Key:  "name",
					Type: "string",
					List: false,
				},
				{
					Id: "",
					Name: map[string]string{
						"zh-CN": "年龄",
					},
					Key:  "age",
					Type: "int",
					List: false,
				},
			},
			InputParams: nil,
			OutputAnchors: []model.NodeAnchor{
				{
					Id: "",
					Name: map[string]string{
						"zh-CN": "信息",
					},
					Key:  "msg",
					Type: "string",
					List: false,
				},
			},
			Selected: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	//t.Logf("%+v", m)

	r, err := m.Cmder.Exec(nil, map[string]interface{}{
		"name": "bysir",
		"age":  18,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, out := range m.Schema.Outputs {
		t.Logf("%v: %v", out.Key, r[out.Key])
	}
}
