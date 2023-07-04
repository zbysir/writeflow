package llm

import (
	"context"
	"github.com/zbysir/writeflow/pkg/export"
	"io"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	cmd := NewLangChain(nil).Cmd()
	lc := cmd["langchain_call"]
	ne := cmd["new_openai"]

	ctx := context.Background()
	rsp, err := ne.Exec(ctx, map[string]interface{}{
		"api_key": os.Getenv("APIKEY"),
	})
	if err != nil {
		t.Fatal(err)
	}

	rsp, err = lc.Exec(ctx, map[string]interface{}{
		"llm":    rsp["default"],
		"stream": true,
		"prompt": "Hello, my name is John. I am a doctor.",
	})
	if err != nil {
		t.Fatal(err)
	}

	switch r := rsp["default"].(type) {
	case string:
		t.Logf("%s", r)
	case export.Stream:
		t.Log("stream")
		re := r.NewReader()
		for {
			s, err := re.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}

			t.Logf("%+v", s)
		}
	default:
		t.Fatalf("unknown type %s", r)
	}
}
