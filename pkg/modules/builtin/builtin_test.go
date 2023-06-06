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
	wf := writeflow.NewWriteFlow(writeflow.WithModules(New()))

	f, err := writeflow.FlowFromModel(&model.Flow{
		Id:          0,
		Name:        "demo_flow",
		Description: "",
		Graph: model.Graph{
			Nodes: []model.Node{
				{
					Id:   "hello_1",
					Type: "_",
					Data: model.NodeData{
						ComponentData: model.ComponentData{
							Source: model.ComponentSource{
								CmdType:    model.BuiltInCmd,
								BuiltinCmd: "raw",
							},
							InputParams: []model.NodeInputParam{
								{
									Key:      "name",
									Type:     "string",
									Optional: false,
								},
								{
									Key:      "age",
									Type:     "int",
									Optional: false,
								},
							},
							OutputAnchors: []model.NodeAnchor{
								{
									Key:  "default",
									Type: "string",
								},
							},
						},
						Inputs: map[string]string{"name": "bysir", "age": "18"},
					},
					PositionAbsolute: model.NodePosition{},
				},
				{
					Id:   "select_1",
					Type: "_",
					Data: model.NodeData{
						ComponentData: model.ComponentData{
							Source: model.ComponentSource{
								CmdType:    model.BuiltInCmd,
								BuiltinCmd: "select",
							},
							InputParams: []model.NodeInputParam{
								{
									Key:      "path",
									Type:     "string",
									Optional: false,
								},
							},
							InputAnchors: []model.NodeAnchor{
								{
									Key:  "data",
									Type: "string",
								},
							},
						},
						Inputs: map[string]string{"path": "name", "data": "hello_1.default"},
					},
					PositionAbsolute: model.NodePosition{},
				},
			},
			OutputNodeId: "select_1",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	r, err := wf.ExecFlow(ctx, f, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bysir", r["default"])
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
