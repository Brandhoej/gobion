package scfg

import (
	"bytes"
	"go/ast"
	"go/token"
	"testing"

	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
)

func Test(t *testing.T) {
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

	cfg := cfg.New[ast.Stmt, ast.Expr](forLoop.Init)

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

	scfg := New(cfg)
	scfg.Into(scfg.Global(), cfg.Entry(), post)
	scfg.Into(scfg.ZoomIn(scfg.Global()), body)
	

	// Act
	var buffer bytes.Buffer
	scfg.DOT(&buffer)
	t.Log(buffer.String())

	// Assert
	t.FailNow()
}
