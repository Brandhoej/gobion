package golang

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

func NewExpressionsInterpreter(context *z3.Context) *GoExpressionInterpreter {
	return &GoExpressionInterpreter{
		context: context,
	}
}

func (interpreter *GoExpressionInterpreter) Interpret(
	scope symbolic.Scope, expression ast.Expr, reduced bool,
) *z3.AST {
	if expression == nil {
		return nil
	}
	return interpreter.Expression(scope, reduced, expression)
}

func (interpreter *GoExpressionInterpreter) Expression(scope symbolic.Scope, valuation bool, expression ast.Expr) *z3.AST {
	switch cast := any(expression).(type) {
	case *ast.BinaryExpr:
		return interpreter.Binary(scope, valuation, cast)
	case *ast.UnaryExpr:
		return interpreter.Unary(scope, valuation, cast)
	case *ast.ParenExpr:
		return interpreter.Parenthesis(scope, valuation, cast)
	case *ast.BasicLit:
		return interpreter.Literal(scope, cast)
	case *ast.Ident:
		return interpreter.Identifier(scope, valuation, cast)
	}
	panic("Unsupported expression type")
}

func (interpreter *GoExpressionInterpreter) Binary(scope symbolic.Scope, valuation bool, binary *ast.BinaryExpr) *z3.AST {
	lhs := interpreter.Expression(scope, valuation, binary.X)
	rhs := interpreter.Expression(scope, valuation, binary.Y)
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
	case token.ADD: // +
		return z3.Add(lhs, rhs)
	case token.SUB: // -
		return z3.Subtract(lhs, rhs)
	case token.MUL: // *
		return z3.Multiply(lhs, rhs)
	case token.QUO: // /
		return z3.Divide(lhs, rhs)
	}
	panic("Unsupported binary operator")
}

func (interpreter *GoExpressionInterpreter) Parenthesis(scope symbolic.Scope, valuation bool, parenthesis *ast.ParenExpr) *z3.AST {
	return interpreter.Expression(scope, valuation, parenthesis.X)
}

func (interpreter *GoExpressionInterpreter) Unary(scope symbolic.Scope, valuation bool, unary *ast.UnaryExpr) *z3.AST {
	operand := interpreter.Expression(scope, valuation, unary.X)
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

func (interpreter *GoExpressionInterpreter) Identifier(scope symbolic.Scope, valuation bool, identifier *ast.Ident) *z3.AST {
	switch identifier.Name {
	case "true":
		return interpreter.context.NewTrue()
	case "false":
		return interpreter.context.NewFalse()
	default:
		symbol, encasing := scope.Lookup(identifier.Name)
		if encasing == nil {
			panic("There is no symbol for the identifier")
		}
		if valuation {
			return scope.Valuation(symbol)
		} else {
			return scope.Variable(symbol)
		}
	}
}
