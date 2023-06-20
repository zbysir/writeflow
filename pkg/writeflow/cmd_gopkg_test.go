package writeflow

import (
	"context"
	"testing"
)

func TestGoCMd(t *testing.T) {
	c, err := NewGoPkg(nil, ".././_pkg", "examplegocmd")
	if err != nil {
		t.Fatal(err)
	}
	rr, err := c.Exec(context.Background(), NewMap(map[string]interface{}{"_args": []interface{}{"a", "b"}}))
	if err != nil {
		t.Fatal(err)
	}

	//cs := c.Schema()
	//
	//for _, o := range cs.Outputs {
	//	t.Logf("%+v(%s): %v", o.Type, o.Type, rr[o.Type])
	//}
	t.Logf("%+v", rr)
	//t.Logf("%+v", cs)
}
