package writeflow

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/pkg/schema"
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

		f.RegisterCmd(FunCMD(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
		}).SetSchema(schema.CMDSchema{
			Key: "hello",
		}))

		f.RegisterCmd(FunCMD(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": strings.Join(params["args"].([]string), " ")}, nil
		}).SetSchema(schema.CMDSchema{
			Key: "append",
		}))

		rsp, err := f.ExecFlow(context.Background(), `
version: 1

flow:
  append:
    inputs:
     args: INPUT[_args]

  hello-1:
    cmd: hello
    inputs:
      name: hello-2

  hello-2:
    cmd: hello
    inputs:
      name: append

  END:
    inputs:
      default: hello-1
`, map[string]interface{}{"_args": []string{"zhang", "liang"}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: hello: zhang liang", rsp["default"])
	})
}
