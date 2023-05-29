package easyfs

import (
	"encoding/json"
	"os"
	"testing"
)

func TestFileTree(t *testing.T) {
	ft, err := GetFileTree(os.DirFS("./"), ".", 4)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(ft, " ", " ")
	t.Logf("%s", bs)
}
