package langchain

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

func TestName(t *testing.T) {
	lc := NewLangChain().Cmd()["langchain_call"]

	rsp, err := lc.Exec(context.Background(), map[string]interface{}{
		"llm":    openai.NewClient("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"),
		"prompt": "Hello, my name is John. I am a doctor.",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", rsp)
}
