package symbolic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestXxx(t *testing.T) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "Test", `
	package asd

	func Foo(a bool) bool {
		return true && false
	}
	`, parser.ParseComments)
	ast.Print(fset, node)
}
