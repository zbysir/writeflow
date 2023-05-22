package writeflow

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"strings"
	"testing"
	"time"
)

func TestXFlow(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		f := NewWriteFlow()

		f.RegisterComponent(NewComponent(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
		}), cmd.Schema{
			Key: "hello",
		}))

		f.RegisterComponent(NewComponent(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": strings.Join(params["args"].([]string), " ")}, nil
		}), cmd.Schema{
			Key: "append",
		}))

		rsp, err := f.ExecFlow(context.Background(), &Flow{
			Nodes: map[string]Node{
				"hello1": {
					Id:           "hello1",
					ComponentKey: "hello",
					Inputs: []NodeInput{
						{
							Key:         "name",
							Type:        "anchor",
							Literal:     "",
							NodeId:      "hello2",
							ResponseKey: "default",
						},
					},
				},
				"hello2": {
					Id:           "hello2",
					ComponentKey: "hello",
					Inputs: []NodeInput{
						{
							Key:         "name",
							Type:        "anchor",
							Literal:     "",
							NodeId:      "append",
							ResponseKey: "default",
						},
					},
				},
				"append": {
					Id:           "append",
					ComponentKey: "append",
					Inputs: []NodeInput{
						{
							Key:         "args",
							Type:        "anchor",
							Literal:     "",
							NodeId:      "INPUT",
							ResponseKey: "_args",
						},
					},
				},
				"END": {
					Id:           "END",
					ComponentKey: "",
					Inputs: []NodeInput{
						{
							Key:         "default",
							Type:        "anchor",
							Literal:     "",
							NodeId:      "hello1",
							ResponseKey: "default",
						},
					},
				},
			},
		}, "END", map[string]interface{}{"_args": []string{"zhang", "liang"}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: hello: zhang liang", rsp["default"])
	})

	t.Run("GetCMDs", func(t *testing.T) {
		f := NewWriteFlow()

		f.RegisterComponent(NewComponent(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
		}), cmd.Schema{
			Inputs: []cmd.SchemaParams{
				{
					Key:         "name",
					Type:        "string",
					NameLocales: nil,
					DescLocales: nil,
				},
			},
			Outputs: []cmd.SchemaParams{
				{
					Key:         "default",
					Type:        "string",
					NameLocales: nil,
					DescLocales: nil,
				},
			},
			Key:         "hello",
			NameLocales: map[string]string{"zh": "Say Hello"},
			DescLocales: map[string]string{"zh": "Append 'hello ' to name"},
		}))

		cmds, _ := f.GetCMDs(context.Background(), nil)
		bs, _ := json.Marshal(cmds)
		t.Logf("%s", bs)
	})
}

func TestFromModelFlow(t *testing.T) {
	f, err := FromModelFlow(&model.Flow{
		Id:          0,
		Name:        "demo_flow",
		Description: "",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Width:  0,
					Height: 0,
					Id:     "hello",
					Position: struct {
						X int `json:"x"`
						Y int `json:"y"`
					}{},
					Type: "hello_component",
					Data: model.NodeData{
						Label: "",
						//Id:           "hello",
						Name:         "",
						Type:         "",
						Category:     "",
						Icon:         "",
						Description:  "",
						Inputs:       map[string]string{"name": "bysir", "age": "18"},
						Source:       model.NodeSource{},
						InputAnchors: nil,
						InputParams: []model.NodeInputParam{
							{
								Id:       "",
								Name:     nil,
								Key:      "name",
								Type:     "string",
								Optional: false,
							},
							{
								Id:       "",
								Name:     nil,
								Key:      "age",
								Type:     "int",
								Optional: false,
							},
						},
						OutputAnchors: []model.NodeAnchor{
							{
								Id:   "",
								Name: nil,
								Key:  "default",
								Type: "string",
								List: false,
							},
						},
						Selected: false,
					},
					PositionAbsolute: struct {
						X int `json:"x"`
						Y int `json:"y"`
					}{},
					Selected: false,
					Dragging: false,
				},
			},
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		t.Fatal(err)
	}

	//components := f.UsedComponents()
	//t.Logf("components: %+v", components)

	wf := NewWriteFlow()
	wf.RegisterComponent(NewComponent(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
	}), cmd.Schema{
		Key: "hello_component",
	}))

	rsp, err := wf.ExecFlow(context.Background(), f, "hello", map[string]interface{}{"name": "bysir"})
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	assert.Equal(t, "hello: bysir", rspDefault)

}
