package symbolic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
)

func Test1(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	solver := context.NewSolver()
	path := NewGoGlobalPath(context)
	n := context.NewConstant(z3.WithName("n"), context.IntegerSort())
	m := context.NewConstant(z3.WithName("m"), context.IntegerSort())
	one := context.NewInt(1, context.IntegerSort())
	expected := z3.Subtract(z3.Subtract(m, one), m)

	source := `
	package example

	func Max(n, m int) int {
		a := 0
		a += 1
		a += m
		a = a - 1
		a -= 1
		if a < m {
			return a - m
		} else {
			return n + m
		}
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)

	// Act
	result := InterpretFunction(path, function, []*z3.AST{n, m})

	// Assert
	equality := z3.Eq(result[0], expected)
	if model := solver.Prove(equality); model != nil {
		t.Error(equality.String(), "has solution", model.String())
	}
}

func Test2(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	solver := context.NewSolver()
	path := NewGoGlobalPath(context)
	expected := context.NewInt(2, context.IntegerSort())

	source := `
	package example

	func Max() int {
		sum := 0
		for x := 0; x < 2; x++ {
			sum += 1
		}
		return sum
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)

	// Act
	result := InterpretFunction(path, function, []*z3.AST{})

	// Assert
	equality := z3.Eq(result[0], expected)
	if model := solver.Prove(equality); model != nil {
		t.Error(equality.String(), "has solution", model.String())
	}
}
