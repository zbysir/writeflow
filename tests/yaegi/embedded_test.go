package yaegi

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"testing"
)

func TestEmbedded(t *testing.T) {
	i := interp.New(interp.Options{
		GoPath: "./",
		// wrap src/pkgname
		SourcecodeFilesystem: nil,
	})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Eval(fmt.Sprintf(`
type A struct {
	*B[string]
}

type B[T any] struct {
	data T
}

// https://github.com/traefik/yaegi/issues/1571
func main() {
	a := &A{
		B: &B[string]{},
	}

	_ = a
}
`))
	if err != nil {
		t.Fatal(err)
	}
}
