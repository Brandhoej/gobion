package symbolic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/stretchr/testify/assert"
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
		a = !b
		if b && a {
			b = !a
		} else {
			c := false
			if !c {
				b = !c
			}
		}
		ret = b
	}
	`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	function := node.Decls[0].(*ast.FuncDecl)

	config := z3.NewConfig()
	context := z3.NewContext(config)
	interpreter := NewStatementInterpreter(context)
	global := NewGoGlobalPath(context)

	// Act
	path := interpreter.Function(global, function)
	fmt.Println(path.String())

	// Assert
	t.FailNow()
}

func TestDsa(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	variable := context.NewConstant(
		z3.WithName("foo"), context.BooleanSort(),
	)

	// Act
	identifier := variable.String()

	// Assert
	assert.Equal(t, "foo", identifier)
	assert.Equal(t, z3.KindBoolean, variable.Sort().Kind())
}
