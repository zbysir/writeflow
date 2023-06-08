package builtin

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	r, err := New().Cmd()["select"].Exec(context.Background(), map[string]interface{}{
		"data": map[string]string{
			"a": "b",
		},
		"path": "data.a",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "b", r["default"])
}

func TestTemplateText(t *testing.T) {
	r, err := New().Cmd()["template_text"].Exec(context.Background(), map[string]interface{}{
		"template": "hello {{name}}",
		"name":     "bysir",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello bysir", r["default"])
}

func TestCallHttp(t *testing.T) {
	r, err := New().Cmd()["call_http"].Exec(context.Background(), map[string]interface{}{
		"url":    "https://writeflow.bysir.top/api/flow",
		"body":   "",
		"method": "GET",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, nil, r["default"])
}

func TestSwitch(t *testing.T) {
	r, err := New().Cmd()["switch"].Exec(context.Background(), map[string]interface{}{
		"data": map[string]string{"a": "b"},
		"aaa":  "data.a",
		"bbb":  "data.b",
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, r["aaa"])
	assert.Equal(t, nil, r["bbb"])
}
