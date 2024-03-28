package cfg

import (
	"bytes"
	"go/ast"
	"go/token"
	"testing"
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

	cfg := New[ast.Stmt, ast.Expr](forLoop.Init)

	condition := cfg.JumpTo(cfg.Entry(), cfg.NewCondition(forLoop.Cond))

	body, _ := cfg.NewBlock(forLoop.Body.List...)
	cfg.IfThenElse(
		condition,
		cfg.NewConditionalJump(ast.NewIdent("true"), body),
		cfg.NewConditionalJump(ast.NewIdent("false"), cfg.Exit()),
	)

	post, _ := cfg.NewBlock(forLoop.Post)
	cfg.Sequence(body, post)

	cfg.JumpTo(post, condition)

	// Act
	var buffer bytes.Buffer
	cfg.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
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

	cfg := New[ast.Stmt, ast.Expr](switchStmt.Init)

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
	cond0 := cfg.NewCondition(check0)

	check1 := &ast.BinaryExpr{
		X:  switchStmt.Tag,
		Op: token.EQL,
		Y:  case1.List[0],
	}
	cond1 := cfg.NewCondition(check1)

	check2 := &ast.BinaryExpr{
		X:  switchStmt.Tag,
		Op: token.EQL,
		Y:  case2.List[0],
	}
	cond2 := cfg.NewCondition(check2)

	cfg.JumpTo(cfg.entry, cond0)

	b0, _ := cfg.NewBlock() // Just a fallthrough in the body
	i0, _ := cfg.NewBlock()
	cfg.IfThenElse(
		cond0,
		cfg.NewConditionalJump(ast.NewIdent("true"), b0),
		cfg.NewConditionalJump(ast.NewIdent("false"), i0),
	)

	cfg.JumpTo(i0, cond1)
	b1, _ := cfg.NewBlock() // Just a labeled fallthrough in the body
	i1, _ := cfg.NewBlock()
	cfg.IfThenElse(
		cond1,
		cfg.NewConditionalJump(ast.NewIdent("true"), b1),
		cfg.NewConditionalJump(ast.NewIdent("false"), i1),
	)

	// b0: fallthrough
	cfg.Sequence(b0, b1)

	cfg.JumpTo(i1, cond2)
	b2, cb2 := cfg.NewBlock(case2.Body...)
	i2, ci2 := cfg.NewBlock([]ast.Stmt{case3.Body[0].(*ast.LabeledStmt).Stmt}...)
	cfg.IfThenElse(
		cond2,
		cfg.NewConditionalJump(ast.NewIdent("true"), b2),
		cfg.NewConditionalJump(ast.NewIdent("false"), i2),
	)

	// b1: fallthrough label
	cfg.Sequence(b1, i2)

	cfg.JumpFrom(cb2, cfg.NewUnconditionalJump(cfg.Exit()))
	cfg.JumpFrom(ci2, cfg.NewUnconditionalJump(cfg.Exit()))

	// Act
	var buffer bytes.Buffer
	cfg.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}