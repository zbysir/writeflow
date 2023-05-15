package explore

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExecFunc(t *testing.T) {

	t.Run("base", func(t *testing.T) {
		rsp, err := execFunc(context.Background(), func(ctx context.Context, name string) string {
			return "hello: " + name
		}, []interface{}{"world"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: world", rsp[0])
	})

	t.Run("array", func(t *testing.T) {
		rsp, err := execFunc(context.Background(), func(ctx context.Context, ns []string) string {
			return "hello: " + strings.Join(ns, " ")
		}, []interface{}{[]interface{}{"a", "b"}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: a b", rsp[0])
	})
	t.Run("obj", func(t *testing.T) {
		rsp, err := execFunc(context.Background(), func(ctx context.Context, s struct {
			First string
			End   string
		}) string {
			return "hello: " + s.First + " " + s.End
		}, []interface{}{map[string]interface{}{
			"First": "a",
			"End":   "b",
		}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: a b", rsp[0])
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
    - _args

  hello-1:
    cmd: hello
    inputs:
    - hello-2[0]

  hello-2:
    cmd: hello
    inputs:
    - append[0]

  END:
    inputs:
      - hello-1[0]
`, []interface{}{"zhang", "liang"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: hello: zhang liang", rsp[0])
	})

	t.Run("custom_struct", func(t *testing.T) {
		f := NewShelFlow()

		f.RegisterCmd("hello", FunCMD(func(ctx context.Context, name UserStruct) string {
			return fmt.Sprint("hello: ", name.Name, " age: ", name.Age)
		}))

		f.RegisterCmd("new", FunCMD(func(ctx context.Context, name string, arg int) UserStruct {
			return UserStruct{
				Name: name,
				Age:  arg,
			}
		}))

		rsp, err := f.ExecFlow(context.Background(), `
version: 1

flow:
  new:
    cmd: new
    inputs:
      - _args[0]
      - _args[1]

  hello:
    cmd: hello
    inputs:
      - new[0]
    depends:
      - new

  END:
    cmd: END
    inputs:
      - hello[0]
`, []interface{}{"bysir", 18})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "hello: bysir age: 18", rsp[0])
	})
}

type UserStruct struct {
	Name string
	Age  int
}
