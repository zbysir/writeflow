package writeflow

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestXFlow(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		f := NewShelFlow()

		f.RegisterCmd("hello", FunCMD(func(ctx context.Context, name string) string {
			return "hello: " + name
		}))

		f.RegisterCmd("append", FunCMD(func(ctx context.Context, args []string) string {
			return strings.Join(args, " ")
		}))

		rsp, err := f.ExecFlow(context.Background(), `
version: 1

flow:
  append:
    inputs:
     0: _args[0]

  hello-1:
    cmd: hello
    inputs:
      0: hello-2[0]

  hello-2:
    cmd: hello
    inputs:
      0: append[0]

  END:
    inputs:
      0: hello-1[0]
`, map[string]interface{}{"_args": []string{"zhang", "liang"}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: hello: zhang liang", rsp["0"])
	})
}
