// Code generated by 'yaegi extract go/importer'. DO NOT EDIT.

//go:build go1.20
// +build go1.20

package stdlib

import (
	"go/importer"
	"reflect"
)

func init() {
	Symbols["go/importer/importer"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Default":     reflect.ValueOf(importer.Default),
		"For":         reflect.ValueOf(importer.For),
		"ForCompiler": reflect.ValueOf(importer.ForCompiler),

		// type definitions
		"Lookup": reflect.ValueOf((*importer.Lookup)(nil)),
	}
}
