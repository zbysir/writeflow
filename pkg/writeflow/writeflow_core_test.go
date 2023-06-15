package writeflow

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/internal/cmd"
	"github.com/zbysir/writeflow/pkg/schema"
	"testing"
)

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
						Anchors: []NodeAnchorTarget{
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
						Anchors: []NodeAnchorTarget{
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
				Cmd: "nothing",
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
				Cmd: "nothing",
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
	rsp, err := r.ExecNode(context.Background(), "a", false, func(result NodeStatusLog) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "b", rsp["default"])
	assert.Equal(t, "data=='b'", rsp["branch"])

	f.Nodes["a"].Inputs[0].Literal = "d"
	r = newRunner(nil, &f, 0)
	rsp, err = r.ExecNode(context.Background(), "a", false, func(result NodeStatusLog) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, rsp["default"])
	assert.Equal(t, "", rsp["branch"])

	f.Nodes["a"].Inputs[0].Literal = "c"

	r = newRunner(nil, &f, 0)
	rsp, err = r.ExecNode(context.Background(), "a", false, func(result NodeStatusLog) {
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
						Anchors: []NodeAnchorTarget{
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
						Anchors: []NodeAnchorTarget{
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
	rsp, err := r.ExecNode(context.Background(), "a", false, func(result NodeStatusLog) {
		//t.Logf("%+v", result)
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []interface{}{"hi: a", "hi: b", "hi: c"}, rsp["default"])
}
