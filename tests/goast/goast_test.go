package goast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestGoAst(t *testing.T) {
	// src is the input for which we want to print the AST.
	src := `
package main

type B[T any] struct {
	data T
}

func (b *B[T]) SetData(data T) {
	b.clearData()
}

func (b *B[T]) clearData() {
	var t T
	b.data = t
}

`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}
