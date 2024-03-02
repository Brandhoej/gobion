package symbolic

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/internal/z3"
)

type ExpressionInterpreter struct {
	context *z3.Context
}

func (interpreter *ExpressionInterpreter) Expression(scope Scope[*z3.AST], expression ast.Expr) *z3.AST {
	switch cast := any(expression).(type) {
	case *ast.BinaryExpr:
		return interpreter.Binary(scope, cast)
	case *ast.UnaryExpr:
		return interpreter.Unary(scope, cast)
	case *ast.ParenExpr:
		return interpreter.Parenthesis(scope, cast)
	case *ast.Ident:
		return interpreter.Identifier(scope, cast)
	}
	panic("Unsupported expression type")
}

func (interpreter *ExpressionInterpreter) Binary(scope Scope[*z3.AST], binary *ast.BinaryExpr) *z3.AST {
	lhs := interpreter.Expression(scope, binary.X)
	rhs := interpreter.Expression(scope, binary.Y)
	switch binary.Op {
	case token.LAND:
		return z3.And(lhs, rhs)
	case token.LOR:
		return z3.Or(lhs, rhs)
	case token.XOR:
		return z3.Xor(lhs, rhs)
	}
	panic("Unsupported binary operator")
}

func (interpreter *ExpressionInterpreter) Parenthesis(scope Scope[*z3.AST], parenthesis *ast.ParenExpr) *z3.AST {
	return interpreter.Expression(scope, parenthesis.X)
}

func (interpreter *ExpressionInterpreter) Unary(scope Scope[*z3.AST], unary *ast.UnaryExpr) *z3.AST {
	operand := interpreter.Expression(scope, unary.X)
	switch unary.Op {
	case token.NOT:
		return z3.Not(operand)
	}
	panic("Unsupported unary operator")
}

func (interpreter *ExpressionInterpreter) Identifier(scope Scope[*z3.AST], identifier *ast.Ident) *z3.AST {
	switch identifier.Name {
	case "true":
		return interpreter.context.NewTrue()
	case "false":
		return interpreter.context.NewFalse()
	default:
		symbol, exists := scope.Symbol(identifier.Name)
		if !exists {
			panic("Identifier does not have a symbol")
		}

		// Interpretation of global variables can be reduced to the variable declaration.
		if scope.IsGlobal() && scope.IsLocal(identifier.Name) {
			variable, exists := scope.Variable(symbol)
			if !exists {
				panic("Symbol does not have a variable declared")
			}
			return variable
		}
		
		// If the variable is not global then we want to reduce it.
		valuation, exists := scope.Valuation(symbol)
		if !exists {
			panic("Symbol does not have a valuation")
		}
		return valuation
	}
}
