package writeflow

import (
	"context"
	"testing"
)

func TestGoCMd(t *testing.T) {
	c, err := NewGoPkgCMD(nil, "./_pkg", "examplegocmd", "examplegocmd.New")
	if err != nil {
		t.Fatal(err)
	}
	rr, err := c.Exec(context.Background(), map[string]interface{}{"__args": []interface{}{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}

	cs := c.Schema(context.Background())
	t.Logf("%+v", rr)
	t.Logf("%+v", cs)
}
