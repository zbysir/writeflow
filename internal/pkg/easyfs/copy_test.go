package easyfs

import (
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/zbysir/writeflow/internal/pkg/db"
	"github.com/zbysir/writeflow/internal/pkg/gobilly"
	"testing"
)

func TestCopy(t *testing.T) {
	d, err := db.NewKvDb("./editor/database")
	if err != nil {
		t.Fatal(err)
	}
	st, err := d.Open(fmt.Sprintf("project_1"), "theme")
	if err != nil {
		t.Fatal(err)
	}
	fsTheme := gobilly.NewDbFs(st)
	if err != nil {
		t.Fatal(err)
	}

	err = CopyDir("/", "", gobilly.NewStdFs(fsTheme), osfs.New("../.cached"))
	if err != nil {
		t.Fatal(err)
	}
}
