package cmd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecFunc(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		rsp, err := execFunc(context.Background(), func(ctx context.Context, name string) string {
			return "hello: " + name
		}, map[string]interface{}{"0": "world"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: world", rsp["0"])
	})

	t.Run("obj", func(t *testing.T) {
		rsp, err := execFunc(context.Background(), func(ctx context.Context, s struct {
			First string
			End   string
		}) string {
			return "hello: " + s.First + " " + s.End
		}, map[string]interface{}{"0": map[string]interface{}{
			"First": "a",
			"End":   "b",
		}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: a b", rsp["0"])
	})
}
