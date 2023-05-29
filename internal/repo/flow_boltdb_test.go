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
		Key:      "openai",
		Category: "",
		Data: model.ComponentData{
			Name:        nil,
			Icon:        "",
			Description: nil,
			Source: model.ComponentSource{
				Type:    "",
				CmdType: "",
				GitUrl:  "",
				GoScript: model.ComponentGoScript{
					Script: "package main;",
				},
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

	assert.Equal(t, "package main;", cs["openai"].Data.Source.GoScript.Script)

	cl, total, err := f.GetComponentList(ctx, GetFlowListParams{})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cs: %v, total: %+v", cl, total)
}
