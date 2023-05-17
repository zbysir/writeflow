package writeflow

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zbysir/writeflow/cmd"
	"github.com/zbysir/writeflow/pkg/schema"
	"strings"
	"testing"
)

func TestXFlow(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		f := NewShelFlow()

		f.RegisterCmd(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{"default": "hello: " + (params["name"].(string))}, nil
		}).SetSchema(schema.CMDSchema{
			Key: "hello",
		}))

		f.RegisterCmd(cmd.NewFun(func(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
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
