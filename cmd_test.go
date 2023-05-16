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
	rr, err := c.Exec(context.Background(), []interface{}{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", rr)
}
