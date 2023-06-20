package writeflow

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJavascript(t *testing.T) {
	s, _ := NewJavaScriptCMD("function exec(params) {return {default:params.a}}")
	ctx := context.Background()
	r, err := s.Exec(ctx, map[string]interface{}{"a": 1})
	if err != nil {
		t.Fatal()
	}

	assert.Equal(t, int64(1), r["default"])
}
