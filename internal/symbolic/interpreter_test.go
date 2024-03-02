package symbolic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
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

func TestAsd(t *testing.T) {
	// Arrange
	source := `
	package example

	func Foo(a, b bool) (ret bool) {
		if b {
			b = !a
		}
		return !a
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)

	config := z3.NewConfig()
	context := z3.NewContext(config)
	NewFunctionInterpreter(context, function)

	// Act

	// Assert
	t.FailNow()
}
