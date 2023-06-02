package writeflow

import (
	"context"
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

		f.RegisterCmd("hello", cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
		}))

		f.RegisterCmd("append", cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": strings.Join(params["args"].([]string), " ")}, nil
		}))

		rsp, err := f.ExecFlow(context.Background(), &Flow{
			OutputNodeId: "END",
			Nodes: map[string]Node{
				"hello1": {
					Id:  "hello1",
					Cmd: "hello",
					Inputs: []NodeInput{
						{
							Key:       "name",
							Type:      "anchor",
							Literal:   "",
							NodeId:    "hello2",
							OutputKey: "default",
						},
					},
				},
				"hello2": {
					Id:  "hello2",
					Cmd: "hello",
					Inputs: []NodeInput{
						{
							Key:       "name",
							Type:      "anchor",
							Literal:   "",
							NodeId:    "append",
							OutputKey: "default",
						},
					},
				},
				"append": {
					Id:  "append",
					Cmd: "append",
					Inputs: []NodeInput{
						{
							Key:       "args",
							Type:      "anchor",
							Literal:   "",
							NodeId:    "INPUT",
							OutputKey: "_args",
						},
					},
				},
				"END": {
					Id:  "END",
					Cmd: "",
					Inputs: []NodeInput{
						{
							Key:       "default",
							Type:      "anchor",
							Literal:   "",
							NodeId:    "hello1",
							OutputKey: "default",
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
}

func TestFromModelFlow(t *testing.T) {
	f, err := FlowFromModel(&model.Flow{
		Id:          0,
		Name:        "demo_flow",
		Description: "",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Width:    0,
					Height:   0,
					Id:       "hello",
					Position: model.NodePosition{},
					Type:     "hello_component",
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
					PositionAbsolute: model.NodePosition{},
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
	wf.RegisterCmd("hello_component", cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
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
	wf.RegisterCmd("call", callLLM)

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
	wf.RegisterCmd("openai", newScript)

	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"query": "特斯拉是谁"})
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	t.Logf("%v", rspDefault)
}
