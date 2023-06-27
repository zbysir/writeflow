package gobilly

import (
	"github.com/zbysir/writeflow/internal/pkg/db"
	"io/ioutil"
	"testing"
)

func TestGoBilly(t *testing.T) {
	d := db.NewKvDb("./database")
	defer d.Clean("test")

	st, err := d.Open("test", "theme")
	if err != nil {
		t.Fatal(err)
	}

	f := NewDbFs(st)

	t.Run("mkdir", func(t *testing.T) {
		err := f.MkdirAll("/", 0)
		if err != nil {
			t.Fatal(err)
		}
		err = f.MkdirAll("component", 0)
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.Create("componentx")
		if err != nil {
			t.Fatal(err)
		}
		fa, err := f.Create("component/a.jsx")
		if err != nil {
			t.Fatal(err)
		}
		fa.Write([]byte(`hello`))
	})

	t.Run("read dir /", func(t *testing.T) {
		ds, err := f.ReadDir("/")
		if err != nil {
			t.Fatal(err)
		}
		for _, i := range ds {
			t.Logf("%v %v", i.IsDir(), i.Name())
		}
	})

	t.Run("read dir component", func(t *testing.T) {
		ds, err := f.ReadDir("component")
		if err != nil {
			t.Fatal(err)
		}
		for _, i := range ds {
			t.Logf("%v %v", i.IsDir(), i.Name())
		}
	})

	t.Run("open file", func(t *testing.T) {
		info, err := f.Stat("./component/a.jsx")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("info %+v", info.ModTime())

		fi, err := f.Open("./component/a.jsx")
		if err != nil {
			t.Fatal(err)
		}

		bs, err := ioutil.ReadAll(fi)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Logf("%s", bs)
	})
}

func TestMakeAll(t *testing.T) {
	d := db.NewKvDb("./database")
	defer d.Clean("makeall")

	st, err := d.Open("makeall", "default")
	if err != nil {
		t.Fatal(err)
	}

	f := NewDbFs(st)

	t.Run("mkdir", func(t *testing.T) {
		err = f.MkdirAll("component", 0)
		if err != nil {
			t.Fatal(err)
		}

		_, err = f.Create("componentx")
		if err != nil {
			t.Fatal(err)
		}
		//err = f.MkdirAll("componentx", 0)
		//if err != nil {
		//	t.Fatal(err)
		//}
		err = f.MkdirAll("component/a/b/c/d", 0)
		if err != nil {
			t.Fatal(err)
		}
		ds, err := f.ReadDir("")
		if err != nil {
			t.Fatal(err)
		}
		for _, i := range ds {
			t.Logf("%v %v", i.IsDir(), i.Name())
		}
	})

}
