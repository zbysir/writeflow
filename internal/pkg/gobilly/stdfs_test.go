package gobilly

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestStdFs(t *testing.T) {
	log.SetDev(true)

	f := osfs.New("../")
	std := NewStdFs(f)

	t.Run("walk", func(t *testing.T) {
		err := fs.WalkDir(std, "", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%v %v %v", d.IsDir(), path, d.Name())
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("walkos", func(t *testing.T) {
		err := filepath.WalkDir("../", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%v %v %v", d.IsDir(), path, d.Name())
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	})

}
