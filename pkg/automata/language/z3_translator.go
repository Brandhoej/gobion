package language

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Z3Translator struct {
	context   *z3.Context
	variables Variables
}

func NewZ3Translator(
	context *z3.Context,
	variables Variables,
) Z3Translator {
	return Z3Translator{
		context:   context,
		variables: variables,
	}
}

func (translator Z3Translator) sort(sort Sort) *z3.Sort {
	if sort == BooleanSort {
		return translator.context.BooleanSort()
	}
	if sort == IntegerSort {
		return translator.context.IntegerSort()
	}
	panic("Unknown sort")
}

func (translator Z3Translator) Translate(expression Expression) *z3.AST {
	switch cast := any(expression).(type) {
	case Variable:
		return translator.variable(cast)
	case Binary:
		return translator.binary(cast)
	case Integer:
		return translator.integer(cast)
	case Boolean:
		return translator.boolean(cast)
	case Unary:
		return translator.unary(cast)
	case IfThenElse:
		return translator.ifThenElse(cast)
	}
	panic("Unknown expression type")
}

func (translator Z3Translator) symbolVariable(symbol symbols.Symbol) *z3.AST {
	if sort, exists := translator.variables.Lookup(symbol); exists {
		return translator.context.NewConstant(z3.WithInt(int(symbol)), translator.sort(sort))
	}
	return nil
}

func (translator Z3Translator) variable(variable Variable) *z3.AST {
	return translator.symbolVariable(variable.symbol)
}

func (translator Z3Translator) binary(binary Binary) *z3.AST {
	lhs := translator.Translate(binary.lhs)
	rhs := translator.Translate(binary.rhs)
	switch binary.operator {
	case Equal:
		return z3.Eq(lhs, rhs)
	case NotEqual:
		return z3.Not(z3.Eq(lhs, rhs))
	case LessThan:
		return z3.LT(lhs, rhs)
	case LessThanEqual:
		return z3.LE(lhs, rhs)
	case GreaterThan:
		return z3.GT(lhs, rhs)
	case GreaterThanEqual:
		return z3.GE(lhs, rhs)
	case Addition:
		return z3.Add(lhs, rhs)
	case Subtraction:
		return z3.Subtract(lhs, rhs)
	case LogicalAnd:
		return z3.And(lhs, rhs)
	case LogicalOr:
		return z3.Or(lhs, rhs)
	case Implication:
		return z3.Implies(lhs, rhs)
	}
	panic("Unknown binary operator")
}

func (translator Z3Translator) integer(integer Integer) *z3.AST {
	return translator.context.NewInt(integer.value, translator.context.IntegerSort())
}

func (translator Z3Translator) boolean(boolean Boolean) *z3.AST {
	if boolean.value {
		return translator.context.NewTrue()
	}
	return translator.context.NewFalse()
}

func (translator Z3Translator) unary(unary Unary) *z3.AST {
	operand := translator.Translate(unary.operand)
	switch unary.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary operator")
}

func (translator Z3Translator) ifThenElse(ifThenElse IfThenElse) *z3.AST {
	return z3.ITE(
		translator.Translate(ifThenElse.condition),
		translator.Translate(ifThenElse.consequence),
		translator.Translate(ifThenElse.alternative),
	)
}
