package scfg

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForLoop(t *testing.T) {
	// Arrange
	forLoop := ast.ForStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("i"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			},
		},
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("i"),
			Op: token.LSS,
			Y: &ast.BasicLit{
				Kind:  token.INT,
				Value: "10",
			},
		},
		Post: &ast.IncDecStmt{
			X:   ast.NewIdent("i"),
			Tok: token.INC,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("foo"),
					},
				},
			},
		},
	}

	builder := NewBuilder(token.NewFileSet())

	init := builder.AddBlock(NewBlock(forLoop.Init))
	body := builder.AddBlock(NewBlock(forLoop.Body.List...))
	post := builder.AddBlock(NewBlock(forLoop.Post))

	condition := builder.Branch(init, NewCondition(forLoop.Cond))

	builder.UnconditionalJump(Initial, init)
	builder.Tautology(init, condition, body)
	builder.Contradiction(init, condition, Terminal)
	builder.UnconditionalJump(body, post)
	builder.Condition(post, condition)

	// Act
	var buffer bytes.Buffer
	builder.Build().DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	assert.Equal(t, 0, init)
	assert.Equal(t, 1, body)
	assert.Equal(t, 2, post)
	assert.Equal(t, 0, condition)
	t.FailNow()
}

func TestSwitch(t *testing.T) {
	// Arrange
	n := ast.NewIdent("n")
	label := ast.NewIdent("label")
	one := &ast.BasicLit{
		Kind:  token.INT,
		Value: "1",
	}
	two := &ast.BasicLit{
		Kind:  token.INT,
		Value: "2",
	}
	three := &ast.BasicLit{
		Kind:  token.INT,
		Value: "3",
	}
	four := &ast.BasicLit{
		Kind:  token.INT,
		Value: "4",
	}
	five := &ast.BasicLit{
		Kind:  token.INT,
		Value: "5",
	}

	case0 := &ast.CaseClause{
		List: []ast.Expr{
			one, two,
		},
		Body: []ast.Stmt{
			&ast.BranchStmt{
				Tok: token.FALLTHROUGH,
			},
		},
	}

	case1 := &ast.CaseClause{
		List: []ast.Expr{
			four,
		},
		Body: []ast.Stmt{
			&ast.BranchStmt{
				Tok:   token.FALLTHROUGH,
				Label: label,
			},
		},
	}

	case2 := &ast.CaseClause{
		List: []ast.Expr{
			five,
		},
	}

	case3 := &ast.CaseClause{
		List: nil,
		Body: []ast.Stmt{
			&ast.LabeledStmt{
				Label: label,
				Stmt: &ast.AssignStmt{
					Lhs: []ast.Expr{n},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{five},
				},
			},
		},
	}

	/*switch n := 3; n {
	  case 1, 2:
	  	fallthrough
	  case 4:
	  	fallthrough label
	  case 5:
	  default:
	  label:
	  	n = 5
	  }*/

	switchStmt := &ast.SwitchStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{n},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{three},
		},
		Tag: n,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				case0,
				case1,
				case2,
				case3,
			},
		},
	}

	builder := NewBuilder(token.NewFileSet())

	init := builder.AddBlock(NewBlock(switchStmt.Init))
	def := builder.AddBlock(NewBlock(case3.Body...))
	builder.UnconditionalJump(Initial, init)

	c0 := builder.AddBlock(NewBlock(case0.Body...))
	check0 := &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  switchStmt.Tag,
			Op: token.EQL,
			Y:  case0.List[0],
		},
		Op: token.LOR,
		Y: &ast.BinaryExpr{
			X:  switchStmt.Tag,
			Op: token.EQL,
			Y:  case0.List[1],
		},
	}
	condition0 := builder.Branch(init, NewCondition(check0))
	builder.Tautology(init, condition0, c0)
	builder.Contradiction(init, condition0, def)

	c1 := builder.AddBlock(NewBlock(case1.Body...))
	check1 := &ast.BinaryExpr{
		X:  switchStmt.Tag,
		Op: token.EQL,
		Y:  case1.List[0],
	}
	condition1 := builder.Branch(init, NewCondition(check1))
	builder.Tautology(init, condition1, c1)
	builder.Contradiction(init, condition1, def)
	builder.UnconditionalJump(c1, def)
	builder.UnconditionalJump(c0, c1)

	c2 := builder.AddBlock(NewBlock(case2.Body...))
	check2 := &ast.BinaryExpr{
		X:  switchStmt.Tag,
		Op: token.EQL,
		Y:  case2.List[0],
	}
	condition2 := builder.Branch(init, NewCondition(check2))
	builder.Tautology(init, condition2, c2)
	builder.Contradiction(init, condition2, def)
	builder.UnconditionalJump(c2, Terminal)

	builder.UnconditionalJump(def, Terminal)

	// Act
	var buffer bytes.Buffer
	builder.Build().DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	assert.Equal(t, 0, init)
	t.FailNow()
}

func TestFunctionGoto(t *testing.T) {
	// Arrange
	source := `
	package example
	
	func Foo() {
		fmt.Println(1)
	End:
		fmt.Println(2)
		fmt.Println(3)
		goto End
		fmt.Println(4)
	}`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	foo := node.Decls[0].(*ast.FuncDecl)

	// Act
	factory := NewFactory(fset)
	scfg := factory.Function(foo)

	// Assert
	var buffer bytes.Buffer
	scfg.DOT(&buffer)
	t.Log(buffer.String())

	t.FailNow()
}

func TestFunctionIf(t *testing.T) {
	// Arrange
	source := `
	package example
	
	func Foo() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	foo := node.Decls[0].(*ast.FuncDecl)

	// Act
	factory := NewFactory(fset)
	scfg := factory.Function(foo)

	// Assert
	var buffer bytes.Buffer
	scfg.DOT(&buffer)
	t.Log(buffer.String())

	t.FailNow()
}

func TestFunctionMultipleTerminals(t *testing.T) {
	// Arrange
	source := `
	package example
	
	func Foo() {
		if a {
			fmt.Println(1)
		} else {
			fmt.Println(2)
		}
	}`
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "example", source, parser.ParseComments)
	foo := node.Decls[0].(*ast.FuncDecl)

	// Act
	factory := NewFactory(fset)
	scfg := factory.Function(foo)

	// Assert
	var buffer bytes.Buffer
	scfg.DOT(&buffer)
	t.Log(buffer.String())

	t.FailNow()
}
