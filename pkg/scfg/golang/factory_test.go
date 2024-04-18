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
	scopes := SCFG(function).CFG()
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
	scopes := SCFG(function).CFG()
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
		} else {
			asd
		}
		print
		return n + m
	}
	`
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	scopes := SCFG(function).CFG()
	var buffer bytes.Buffer
	scopes.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}

func TestSwitch(t *testing.T) {
	// Arrange
	source := `
	package example

	func Max(n, m int) int {
		print(a)
		switch i := 0; i {
		case 1, 2:
			print(b)
			fallthrough
			foo:
				print(c)
		case 3:
			print(d)
			goto foo
		case 4:
			print(e)
			return
			print(f)
		}
		print(g)
	}
	`
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	flow := SCFG(function).CFG()
	var buffer bytes.Buffer
	flow.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}

func TestGoto(t *testing.T) {
	// Arrange
	source := `
	package example

	func Max(n, m int) int {
		print(1)
		goto a
		print(2)
		a:
		print(3)
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

func TestForContinueBreak(t *testing.T) {
	// Arrange
	source := `
	package example

	func Max(n, m int) int {
		print(0)
		for i := 0; i < 10; i++ {
			if i > 5 {
				break;
				print(1)
			}
			print(2)
			continue
			print(3)
		}
		print(4)
	}
	`
	file, _ := parser.ParseFile(token.NewFileSet(), "foo", source, parser.ParseComments)
	function := file.Decls[0].(*ast.FuncDecl)

	// Act
	flow := SCFG(function)
	var buffer bytes.Buffer
	flow.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}
