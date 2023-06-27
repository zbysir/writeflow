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
type A struct {
	*B[string]
}

type B[T any] struct {
	data T
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
