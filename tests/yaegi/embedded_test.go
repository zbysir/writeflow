package yaegi

import (
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

	_, err = i.Eval(`package main
	import "fmt"
	
	type B[T any] struct {
		data T
	}
	
	func (b *B[T]) SetData(data T) {
		b.clearData()
		b.data = data
	}
	
	func (b *B[T]) clearData() {
		var t T
		b.data = t
	}
	
	func main() {
		b :=  &B[string]{}
		b.SetData("123")
		fmt.Printf("data: %s\n", b.data)
		b.clearData()	
		fmt.Printf("data: %s\n", b.data)
	}
	`)
	if err != nil {
		if pe, ok := err.(interp.Panic); ok {
			t.Logf("%s", pe.Stack)
		}
		t.Fatal(err)
	}
}

func TestSt(t *testing.T) {
	i := interp.New(interp.Options{
		GoPath: "./",
		// wrap src/pkgname
		SourcecodeFilesystem: nil,
	})

	err := i.Use(stdlib.Symbols)
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.Eval(`package main
	import "fmt"
	
	type B struct {
		data any
	}
	
	func (b B) SetData(data any) {
		b.clearData()
		b.data = data
	}
	
	func (b B) clearData() {
		var t any
		b.data = t
	}
	
	func main() {
		b :=  &B{}
		b.SetData("123")
		fmt.Printf("data: %s\n", b.data)
		b.clearData()	
		fmt.Printf("data: %s\n", b.data)
	}
	`)
	if err != nil {
		if pe, ok := err.(interp.Panic); ok {
			t.Logf("%s", pe.Stack)
		}
		t.Fatal(err)
	}
}
