package cmd

import (
	"testing"
)

func TestNewGoScript(t *testing.T) {
	cmd, err := NewGoScript(nil, "", `package main
					import (
				
				"context"
				"fmt"
				)
					
					func Exec(ctx context.Context, params map[string]interface{}) (rsp map[string]interface{}, err error) {
						return map[string]interface{}{"msg": fmt.Sprintf("hello %v %v", params["name"], params["age"])}, nil
					}
					`)
	if err != nil {
		t.Fatal(err)
	}

	r, err := cmd.Exec(nil, map[string]interface{}{
		"name": "bysir",
		"age":  18,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", r["msg"])
}
