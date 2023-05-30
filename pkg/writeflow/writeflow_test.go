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
			OutputNodeId: "END",
			Nodes: map[string]Node{
				"hello1": {
					Id:            "hello1",
					ComponentType: "hello",
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
					Id:            "hello2",
					ComponentType: "hello",
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
					Id:            "append",
					ComponentType: "append",
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
					Id:            "END",
					ComponentType: "",
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
		}, map[string]interface{}{"_args": []string{"zhang", "liang"}})
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
	f, err := FlowFromModel(&model.Flow{
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
						ComponentData: model.ComponentData{
							Source:       model.ComponentSource{},
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
						},
						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
					PositionAbsolute: struct {
						X int `json:"x"`
						Y int `json:"y"`
					}{},
				},
			},
			OutputNodeId: "hello",
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

	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"name": "bysir"})
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	assert.Equal(t, "hello: bysir", rspDefault)
}

func TestOpenAIFlow(t *testing.T) {
	f, err := FlowFromModel(&model.Flow{
		Id:          0,
		Name:        "demo_flow",
		Description: "",
		Graph: model.Graph{
			OutputNodeId: "call",
			Nodes: []model.Node{
				{
					Id:   "openai",
					Type: "openai",
					Data: model.NodeData{
						Inputs: map[string]string{"apikey": "sk-xx"},
						ComponentData: model.ComponentData{
							InputParams: []model.NodeInputParam{
								{Key: "apikey", Type: "string", Optional: false},
								{Key: "base_url", Type: "string", Optional: false},
							},
							OutputAnchors: []model.NodeAnchor{
								{Key: "default", Type: "llm", List: false},
							},
						},
					},
				},
				{
					Id:   "call",
					Type: "call",
					Data: model.NodeData{
						Inputs: map[string]string{"query": "INPUT.query", "llm": "openai.default"},
						ComponentData: model.ComponentData{
							InputAnchors: []model.NodeAnchor{
								{Key: "query", Type: "string"},
								{Key: "llm", Type: "llm"},
							},
							OutputAnchors: []model.NodeAnchor{
								{Key: "default", Type: "string", List: false},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	wf := NewWriteFlow()
	// 如果要使用 llms.LLM interface，必须先注册
	callLLM, err := cmd.NewGoScript(nil, "/Users/bysir/goproj/bysir/writeflow/_pkg", `package main
						import (
							"context"
							"fmt"
							"github.com/tmc/langchaingo/llms"
					)
	
						func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
							query:=params["query"].(string)
							l,_:=params["llm"].(llms.LLM)
							rsps,err:=l.Call(ctx, query)
							if err != nil {
								return nil, err
							}
							return map[string]interface{}{"default": rsps}, nil
						}
						`)
	if err != nil {
		t.Fatal(err)
	}
	wf.RegisterComponent(NewComponent(callLLM, cmd.Schema{
		Key: "call",
	}))

	newScript, err := cmd.NewGoScript(nil, "/Users/bysir/goproj/bysir/writeflow/_pkg", `package main
					import (
				
				"context"
				"fmt"
"github.com/tmc/langchaingo/llms/openai"
				)
					func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
						llm, err := openai.New(openai.WithToken(params["apikey"].(string)))
						if err != nil {	return nil, err}
						return map[string]interface{}{"default": llm}, nil
					}
					`)
	if err != nil {
		t.Fatal(err)
	}
	wf.RegisterComponent(NewComponent(newScript, cmd.Schema{
		Key: "openai",
	}))

	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"query": "特斯拉是谁"})
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	t.Logf("%v", rspDefault)
}
