package builtin

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/internal/model"
	"github.com/zbysir/writeflow/pkg/writeflow"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	r, err := New().Cmd()["select"].Exec(context.Background(), map[string]interface{}{
		"data": map[string]string{
			"a": "b",
		},
		"path": "data.a",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "b", r["default"])
}

func TestTemplateText(t *testing.T) {
	r, err := New().Cmd()["template_text"].Exec(context.Background(), map[string]interface{}{
		"template": "hello {{name}}",
		"name":     "bysir",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello bysir", r["default"])
}

func TestCallHttp(t *testing.T) {
	r, err := New().Cmd()["call_http"].Exec(context.Background(), map[string]interface{}{
		"url":    "http://localhost:18002/rpc",
		"body":   "{     \"service\": \"lb.content.datadriver\",     \"endpoint\": \"QueryService.GetArticles\",     \"request\": {         \"conditions\":  [             {                 \"original_type\": 1,                 \"counter_ids\":[\"ST/US/TSLA\"],                 \"statuses\": [2],                 \"kinds\": [1]             }                      ],         \"started_at\":1686227717282478,         \"limit\": 10,         \"with_source\": true     } }",
		"method": "POST",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", r["default"])

	assert.NotEqual(t, nil, r["default"])
}

func TestSwitch(t *testing.T) {
	r, err := New().Cmd()["switch"].Exec(context.Background(), map[string]interface{}{
		"data": map[string]string{"a": "b"},
		"aaa":  "data.a",
		"bbb":  "data.b",
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, r["aaa"])
	assert.Equal(t, nil, r["bbb"])
}

func TestListComponent(t *testing.T) {
	f, err := writeflow.FlowFromModel(&model.Flow{
		Name: "demo_flow",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Id:   "list",
					Type: "list",
					Data: model.NodeData{
						Source: model.ComponentSource{
							CmdType:    model.BuiltInCmd,
							BuiltinCmd: "list",
						},
						InputParams: []model.NodeInputParam{
							{
								Key:       "data",
								Type:      "string",
								List:      true,
								InputType: model.NodeInputTypeAnchor,
								Anchors: []model.NodeAnchorTarget{
									{
										NodeId:    "a",
										OutputKey: "default",
									},
									{
										NodeId:    "b",
										OutputKey: "default",
									},
								},
							},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{
								Name: nil,
								Key:  "default",
								Type: "any",
								List: false,
							},
						},

						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
				},
				{
					Id:   "a",
					Type: "_",
					Data: model.NodeData{
						Source: model.ComponentSource{
							CmdType:    model.BuiltInCmd,
							BuiltinCmd: "raw",
						},
						InputParams: []model.NodeInputParam{
							{
								Key:   "default",
								Type:  "string",
								Value: "a",
							},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{
								Name: nil,
								Key:  "default",
								Type: "any",
								List: false,
							},
						},
					},
				},
				{
					Id:   "b",
					Type: "_",
					Data: model.NodeData{
						Source: model.ComponentSource{
							CmdType:    model.BuiltInCmd,
							BuiltinCmd: "raw",
						},
						InputParams: []model.NodeInputParam{
							{
								Key:   "default",
								Type:  "string",
								Value: "b",
							},
						},
						OutputAnchors: []model.NodeOutputAnchor{
							{
								Name: nil,
								Key:  "default",
								Type: "any",
								List: false,
							},
						},
					},
				},
			},
			OutputNodeId: "list",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		t.Fatal(err)
	}

	wf := writeflow.NewWriteFlowCore()
	for k, v := range New().Cmd() {
		wf.RegisterCmd(k, v)
	}
	rsp, err := wf.ExecFlow(context.Background(), f, nil, 1)
	if err != nil {
		t.Fatal(err)
	}

	rspDefault := rsp["default"]

	assert.Equal(t, []interface{}{"a", "b"}, rspDefault)
}
