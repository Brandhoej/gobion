package symbolic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/stretchr/testify/assert"
)

func TestFunctionDeclaration(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)

	source := `
	package example

	func Max(n, m int) int {
		ret = n + m
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)

	// Act
	ast := InterpretFunction(context, function)

	// Assert
	assert.Equal(t, "(+ (:var 0) (:var 1))", ast.String())
}
