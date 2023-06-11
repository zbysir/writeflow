package writeflow

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/schema"
	"testing"
	"time"
)

func TestGetRootNodes(t *testing.T) {
	f, err := FlowFromModel(&model.Flow{
		Id:          0,
		Name:        "demo_flow",
		Description: "",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Id:   "hello",
					Type: "hello_component",
					Data: model.NodeData{

						InputParams: []model.NodeInputParam{
							{
								Name:     nil,
								Key:      "name",
								Type:     "string",
								Optional: false,
							},
							{
								Name:     nil,
								Key:      "age",
								Type:     "int",
								Optional: false,
							},
						},

						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
				}, {
					Id:   "hello2",
					Type: "hello_component",
					Data: model.NodeData{

						InputParams: []model.NodeInputParam{
							{
								Name:     nil,
								Key:      "name",
								Type:     "string",
								Optional: false,
							},
							{
								Name:     nil,
								Key:      "age",
								Type:     "int",
								Optional: false,
							},
						},

						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
				},
				{
					Id:   "hello3",
					Type: "hello_component",
					Data: model.NodeData{

						InputParams: []model.NodeInputParam{
							{
								Name:      nil,
								Key:       "to_hello2",
								InputType: model.NodeInputTypeAnchor,
							},
						},

						Inputs: map[string]string{"to_hello2": "hello2.name"},
					},
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

	nodes := f.Nodes.GetRootNodes()
	assert.Equal(t, 2, len(nodes))
	assert.Equal(t, "hello", nodes[0].Id)
	assert.Equal(t, "hello3", nodes[1].Id)
}

func TestSwitch(t *testing.T) {
	f := Flow{
		Nodes: map[string]Node{
			"a": {
				Cmd: "_switch",
				Inputs: []NodeInput{
					{
						Key:     "data",
						Type:    "literal",
						Literal: "b",
					},
					{
						Key:     "data=='c'",
						Type:    "anchor",
						Literal: "",
						Anchors: []model.NodeAnchorTarget{
							{
								NodeId:    "c",
								OutputKey: "default",
							},
						},
					},
					{
						Key:     "data=='b'",
						Type:    "anchor",
						Literal: "",
						Anchors: []model.NodeAnchorTarget{
							{
								NodeId:    "b",
								OutputKey: "default",
							},
						},
					},
				},
			},
			"b": {
				Id:  "",
				Cmd: model.NothingCmd,
				Inputs: []NodeInput{
					{
						Key:     "default",
						Type:    "literal",
						Literal: "b",
					},
				},
			},
			"c": {
				Id:  "",
				Cmd: model.NothingCmd,
				Inputs: []NodeInput{
					{
						Key:     "default",
						Type:    "literal",
						Literal: "c",
					},
				},
			},
		},
		OutputNodeId: "a",
	}

	r := newRunner(nil, &f, 1)
	rsp, err := r.ExecNode(context.Background(), "a", false, func(result model.NodeStatus) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "b", rsp["default"])
	assert.Equal(t, "data=='b'", rsp["branch"])

	f.Nodes["a"].Inputs[0].Literal = "d"
	r = newRunner(nil, &f, 0)
	rsp, err = r.ExecNode(context.Background(), "a", false, func(result model.NodeStatus) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, rsp["default"])
	assert.Equal(t, "", rsp["branch"])

	f.Nodes["a"].Inputs[0].Literal = "c"

	r = newRunner(nil, &f, 0)
	rsp, err = r.ExecNode(context.Background(), "a", false, func(result model.NodeStatus) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "c", rsp["default"])
	assert.Equal(t, "data=='c'", rsp["branch"])
}

func TestFor(t *testing.T) {
	f := Flow{
		Nodes: map[string]Node{
			"a": {
				Id:  "a",
				Cmd: "_for",
				Inputs: []NodeInput{
					{
						Key:     "data",
						Type:    "literal",
						Literal: []string{"a", "b", "c"},
					},
					{
						Key:     "item",
						Type:    "anchor",
						Literal: "",
						Anchors: []model.NodeAnchorTarget{
							{
								NodeId:    "b",
								OutputKey: "default",
							},
						},
					},
				},
			},
			"b": {
				Id:  "b",
				Cmd: "add_prefix",
				Inputs: []NodeInput{
					{
						Key:     "prefix",
						Type:    "literal",
						Literal: "hi: ",
					},
					{
						Key:     "default",
						Type:    "anchor",
						Literal: "",
						Anchors: []model.NodeAnchorTarget{
							{
								NodeId:    "a",
								OutputKey: "item",
							},
						},
					},
				},
			},
		},
		OutputNodeId: "a",
	}

	r := newRunner(map[string]schema.CMDer{
		"add_prefix": cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
			return map[string]interface{}{"default": fmt.Sprintf("%v%v", params["prefix"], params["default"])}, nil
		}),
	}, &f, 1)
	rsp, err := r.ExecNode(context.Background(), "a", false, func(result model.NodeStatus) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []interface{}{"hi: a", "hi: b", "hi: c"}, rsp["default"])
}

func TestFromModelFlow(t *testing.T) {
	f, err := FlowFromModel(&model.Flow{
		Name: "demo_flow",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Id:   "hello",
					Type: "hello_component",
					Data: model.NodeData{
						Source: model.ComponentSource{
							CmdType:    model.BuiltInCmd,
							BuiltinCmd: "hello_component",
							GoPackage:  model.ComponentGoPackage{},
							Script:     model.ComponentScript{},
						},
						InputAnchors: nil,
						InputParams: []model.NodeInputParam{
							{
								Key:  "name",
								Type: "string",
							},
							{
								Key:  "age",
								Type: "int",
							},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{
								Name: nil,
								Key:  "default",
								Type: "string",
								List: false,
							},
						},

						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
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

	wf := NewWriteFlowCore()
	wf.RegisterCmd("hello_component", cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
	}))

	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"name": "bysir"}, 1)
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	assert.Equal(t, "hello: bysir", rspDefault)
}

func TestParallel(t *testing.T) {
	f, err := FlowFromModel(&model.Flow{
		Name: "demo_flow",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Id:   "root",
					Type: "sleep",
					Data: model.NodeData{
						Source: model.ComponentSource{
							CmdType:    model.BuiltInCmd,
							BuiltinCmd: "sleep",
						},
						InputAnchors: nil,
						InputParams: []model.NodeInputParam{
							{
								Key:   "second",
								Type:  "int",
								Value: "1",
							},
							{
								Key:       "a",
								InputType: NodeInputAnchor,
								Anchors: []model.NodeAnchorTarget{
									{
										NodeId:    "a",
										OutputKey: "default",
									},
								},
							},
							{
								Key:       "b",
								InputType: NodeInputAnchor,
								Anchors: []model.NodeAnchorTarget{
									{
										NodeId:    "b",
										OutputKey: "default",
									},
								},
							},
							{
								Key:       "c",
								InputType: NodeInputAnchor,
								Anchors: []model.NodeAnchorTarget{
									{
										NodeId:    "c",
										OutputKey: "default",
									},
								},
							},
							{
								Key:       "d",
								InputType: NodeInputAnchor,
								Anchors: []model.NodeAnchorTarget{
									{
										NodeId:    "d",
										OutputKey: "default",
									},
								},
							},
						},
					},
				},
				{
					Id: "a", Type: "sleep",
					Data: model.NodeData{
						Source:      model.ComponentSource{CmdType: model.BuiltInCmd, BuiltinCmd: "sleep"},
						InputParams: []model.NodeInputParam{{Key: "second", Type: "int", Value: "1"}},
					},
				},
				{
					Id: "b", Type: "sleep",
					Data: model.NodeData{
						Source:      model.ComponentSource{CmdType: model.BuiltInCmd, BuiltinCmd: "sleep"},
						InputParams: []model.NodeInputParam{{Key: "second", Type: "int", Value: "1"}},
					},
				},
				{
					Id: "c", Type: "sleep",
					Data: model.NodeData{
						Source:      model.ComponentSource{CmdType: model.BuiltInCmd, BuiltinCmd: "sleep"},
						InputParams: []model.NodeInputParam{{Key: "second", Type: "int", Value: "1"}},
					},
				},
				{
					Id: "d", Type: "sleep",
					Data: model.NodeData{
						Source:      model.ComponentSource{CmdType: model.BuiltInCmd, BuiltinCmd: "sleep"},
						InputParams: []model.NodeInputParam{{Key: "second", Type: "int", Value: "1"}},
					},
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

	wf := NewWriteFlowCore()
	wf.RegisterCmd("sleep", cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
		time.Sleep(time.Second * time.Duration(cast.ToInt(params["second"])))
		return params, nil
	}))

	start := time.Now()
	// 0: [root, a, b,c,d]  5s
	// 1: [root, [a], b,c,d]  4s
	// 2: [root, [a, b], [c, d]] 3s
	// 3: [root, [a, b, c], d] 3s
	// 4: [root, [a, b, c, d]] 2s
	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"name": "bysir"}, 4)
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["name"]

	t.Logf("spend: %s", time.Since(start))
	assert.Equal(t, "bysir", rspDefault)
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

						InputParams: []model.NodeInputParam{
							{Key: "apikey", Type: "string", Optional: false},
							{Key: "base_url", Type: "string", Optional: false},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{Key: "default", Type: "llm", List: false},
						},
					},
				},
				{
					Id:   "call",
					Type: "call",
					Data: model.NodeData{
						Inputs: map[string]string{"query": "INPUT.query", "llm": "openai.default"},

						InputParams: []model.NodeInputParam{
							{Key: "query", Type: "string", InputType: model.NodeInputTypeAnchor},
							{Key: "llm", Type: "llm", InputType: model.NodeInputTypeAnchor},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{Key: "default", Type: "string", List: false},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	wf := NewWriteFlowCore()
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

	rsp, err := wf.ExecFlow(context.Background(), f, map[string]interface{}{"query": "特斯拉是谁"}, 0)
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	t.Logf("%v", rspDefault)
}
