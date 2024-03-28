package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/scfg/golang"
)

func TestSCFG1(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	path := NewGoGlobalPath(context)
	n := context.NewConstant(z3.WithName("n"), context.IntegerSort())
	m := context.NewConstant(z3.WithName("m"), context.IntegerSort())
	path.scope.Declare("n", n)
	path.scope.Declare("m", m)

	source := `
	package example

	func Max(n, m int) (int, int) {
		if n < m {
			return n - m, m
		}
		return n + m, n
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)
	scopes := golang.SCFG(function)

	// Act
	_, outputs := InterpretSCFG(path, scopes, 2)

	// Assert
	t.Log(outputs[0].String())
	t.Log(outputs[1].String())
	t.FailNow()
}

func TestSCFG2(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	path := NewGoGlobalPath(context)

	source := `
	package example

	func Max() int {
		sum := 1
		for x := 0; x < 3000; x++ {
			sum += 1
		}
		return sum
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)
	scopes := golang.SCFG(function)

	// Act
	_, outputs := InterpretSCFG(path, scopes, 1)

	// Assert
	t.Log(outputs[0].String())
	t.FailNow()
}

func BenchmarkSCFGInterpretationOfForLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Arrange
		config := z3.NewConfig()
		context := z3.NewContext(config)
		path := NewGoGlobalPath(context)

		source := `
		package example

		func Max() int {
			sum := 0
			for x := 0; x < 100; x++ {
				sum += 1
			}
			return sum
		}
		`
		fset := token.NewFileSet()
		node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
		function := node.Decls[0].(*ast.FuncDecl)
		scopes := golang.SCFG(function)

		// Act
		InterpretSCFG(path, scopes, 1)
	}
}
