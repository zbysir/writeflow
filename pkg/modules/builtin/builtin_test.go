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
