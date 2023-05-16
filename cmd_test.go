package explore

import (
	"context"
	"testing"
)

func TestGoCMd(t *testing.T) {
	c := GoCMD{}
	rr, err := c.Exec(context.Background(), []interface{}{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", rr)
}
