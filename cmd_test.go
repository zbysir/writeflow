package writeflow

import (
	"context"
	"testing"
)

func TestGoCMd(t *testing.T) {
	c, err := NewGoPkgCMD(nil, "./_pkg", "examplegocmd")
	if err != nil {
		t.Fatal(err)
	}
	rr, err := c.Exec(context.Background(), map[string]interface{}{"_args": []interface{}{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}

	cs := c.Schema()

	for _, o := range cs.Outputs {
		t.Logf("%+v(%s): %v", o.Key, o.Type, rr[o.Key])
	}
	//t.Logf("%+v", rr)
	//t.Logf("%+v", cs)
}
