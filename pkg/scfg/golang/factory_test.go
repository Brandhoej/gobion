package golang

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIf(t *testing.T) {
	// Arrange
	source := `
	package bar

	func Foo() {
		fmt.Println("0")
		if bar := true; true {
			fmt.Println("1")
		} else {
			fmt.Println("2")
		}
		fmt.Println("3")
	}
	`
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	scopes := SCFG(function)
	var buffer bytes.Buffer
	scopes.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}

func TestFor(t *testing.T) {
	// Arrange
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
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	scopes := SCFG(function)
	var buffer bytes.Buffer
	scopes.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}

func TestIfReturn(t *testing.T) {
	// Arrange
	source := `
	package example

	func Max(n, m int) int {
		if n < m {
			return n - m
		}
		return n + m
	}
	`
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	scopes := SCFG(function)
	var buffer bytes.Buffer
	scopes.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}
