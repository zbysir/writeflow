package repo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/internal/model"
	"testing"
)

func TestComponent(t *testing.T) {
	x, err := NewKvDb("testdata")
	if err != nil {
		t.Fatal(err)
	}
	s, err := x.Open("db", "default")
	if err != nil {
		t.Fatal(err)
	}
	f := NewBoltDBFlow(s)
	ctx := context.Background()
	err = f.CreateComponent(ctx, &model.Component{
		Type:     "openai",
		Category: "",
		Data: model.ComponentData{
			Name:        nil,
			Icon:        "",
			Description: nil,
			Source: model.ComponentSource{
				CmdType: "",
				Script:  model.ComponentScript{InputKey: "script"},
			},
			InputAnchors:  nil,
			InputParams:   nil,
			OutputAnchors: nil,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	cs, err := f.GetComponentByKeys(ctx, []string{"openai"})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "script", cs["openai"].Data.Source.Script.InputKey)

	cl, total, err := f.GetComponentList(ctx, GetFlowListParams{})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cs: %v, total: %+v", cl, total)
}

func TestFlow(t *testing.T) {
	x, err := NewKvDb("testdata")
	if err != nil {
		t.Fatal(err)
	}
	s, err := x.Open("db", "default")
	if err != nil {
		t.Fatal(err)
	}
	f := NewBoltDBFlow(s)
	ctx := context.Background()
	_, err = f.CreateFlow(ctx, &model.Flow{
		Id: 0,
		Graph: model.Graph{
			Nodes: nil,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	cl, total, err := f.GetFlowList(ctx, GetFlowListParams{
		Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(cl))

	for _, c := range cl {
		ef, exist, err := f.GetFlowById(ctx, c.Id)
		if err != nil {
			t.Fatal(err)
		}
		if !exist {
			t.Fatalf("id %v not exist", c.Id)
		}

		assert.Equal(t, ef.Id, c.Id)

		err = f.DeleteFlow(ctx, c.Id)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, total, err = f.GetFlowList(ctx, GetFlowListParams{
		Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int(0), total)
}
