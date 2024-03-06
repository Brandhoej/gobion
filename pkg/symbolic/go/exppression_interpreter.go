package symbolic

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbolic"
)

type GoExpressionInterpreter struct {
	context *z3.Context
}

func (interpreter *GoExpressionInterpreter) Expression(scope symbolic.Scope, expression ast.Expr) *z3.AST {
	switch cast := any(expression).(type) {
	case *ast.BinaryExpr:
		return interpreter.Binary(scope, cast)
	case *ast.UnaryExpr:
		return interpreter.Unary(scope, cast)
	case *ast.ParenExpr:
		return interpreter.Parenthesis(scope, cast)
	case *ast.BasicLit:
		return interpreter.Literal(scope, cast)
	case *ast.Ident:
		return interpreter.Identifier(scope, cast)
	}
	panic("Unsupported expression type")
}

func (interpreter *GoExpressionInterpreter) Binary(scope symbolic.Scope, binary *ast.BinaryExpr) *z3.AST {
	lhs := interpreter.Expression(scope, binary.X)
	rhs := interpreter.Expression(scope, binary.Y)
	switch binary.Op {
	// Logical operators on booleans.
	case token.LAND: // &&
		return z3.And(lhs, rhs)
	case token.LOR: // ||
		return z3.Or(lhs, rhs)
	case token.XOR: // ^
		return z3.Xor(lhs, rhs)

	// Relational operators.
	case token.EQL: // ==
		return z3.Eq(lhs, rhs)
	case token.NEQ: // !=
		return z3.Not(z3.Eq(lhs, rhs))
	case token.LEQ: // <=
		return z3.LE(lhs, rhs)
	case token.GEQ: // >=
		return z3.GE(lhs, rhs)
	case token.LSS: // <
		return z3.LT(lhs, rhs)
	case token.GTR: // >
		return z3.GT(lhs, rhs)
	// Arthimitic operators:
	case token.ADD:
		return z3.Add(lhs, rhs)
	case token.SUB:
		return z3.Subtract(lhs, rhs)
	case token.MUL:
		return z3.Multiply(lhs, rhs)
	case token.QUO:
		return z3.Divide(lhs, rhs)
	}
	panic("Unsupported binary operator")
}

func (interpreter *GoExpressionInterpreter) Parenthesis(scope symbolic.Scope, parenthesis *ast.ParenExpr) *z3.AST {
	return interpreter.Expression(scope, parenthesis.X)
}

func (interpreter *GoExpressionInterpreter) Unary(scope symbolic.Scope, unary *ast.UnaryExpr) *z3.AST {
	operand := interpreter.Expression(scope, unary.X)
	switch unary.Op {
	case token.NOT:
		return z3.Not(operand)
	case token.SUB:
		return z3.Minus(operand)
	}
	panic("Unsupported unary operator")
}

func (interpreter *GoExpressionInterpreter) Literal(scope symbolic.Scope, literal *ast.BasicLit) *z3.AST {
	switch literal.Kind {
	case token.INT:
		integer, _ := strconv.Atoi(literal.Value)
		sort := interpreter.context.IntegerSort()
		return interpreter.context.NewInt(integer, sort)
	}
	panic("Unsupported literal type")
}

func (interpreter *GoExpressionInterpreter) Identifier(scope symbolic.Scope, identifier *ast.Ident) *z3.AST {
	switch identifier.Name {
	case "true":
		return interpreter.context.NewTrue()
	case "false":
		return interpreter.context.NewFalse()
	default:
		valuation, exists := scope.Valuation(identifier.Name)
		if !exists {
			panic("Symbol does not have a valuation")
		}
		return valuation
	}
}
