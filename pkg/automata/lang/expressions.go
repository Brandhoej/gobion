package automata

import "github.com/Brandhoej/gobion/pkg/symbols"

type ExpressionVisitor interface {
	Variable(variable Variable)
	Binary(binary Binary)
	Integer(integer Integer)
}

type Expression interface {
	Accept(visitor ExpressionVisitor)
}

type Integer struct {
	value int
}

func NewInteger(value int) Integer {
	return Integer{
		value: value,
	}
}

func (integer Integer) Accept(visitor ExpressionVisitor) {
	visitor.Integer(integer)
}

type Variable struct {
	symbol symbols.Symbol
}

func NewVariable(symbol symbols.Symbol) Variable {
	return Variable{
		symbol: symbol,
	}
}

func (variable Variable) Accept(visitor ExpressionVisitor) {
	visitor.Variable(variable)
}

type BinaryOperator uint16

const (
	Equal            = BinaryOperator(0)
	NotEqual         = BinaryOperator(1)
	LessThan         = BinaryOperator(2)
	LessThanEqual    = BinaryOperator(3)
	GreaterThan      = BinaryOperator(4)
	GreaterThanEqual = BinaryOperator(5)
	LogicalAnd       = BinaryOperator(6)
	LogicalOr        = BinaryOperator(6)
)

type Binary struct {
	lhs, rhs Expression
	operator BinaryOperator
}

func NewBinary(lhs Expression, operator BinaryOperator, rhs Expression) Binary {
	return Binary{
		lhs: lhs,
		operator: operator,
		rhs: rhs,
	}
}

func NewConjunctions(expression Expression, expressions ...Expression) (conjunction Expression) {
	conjunction = expression
	if len(expressions) == 0 {
		return conjunction
	}

	for idx := range expressions {
		conjunction = NewBinary(conjunction, LogicalAnd, expressions[idx])
	}
	return conjunction
}

func NewDisjunctions(expression Expression, expressions ...Expression) (disjunction Expression) {
	disjunction = expression
	if len(expressions) == 0 {
		return disjunction
	}

	for idx := range expressions {
		disjunction = NewBinary(disjunction, LogicalAnd, expressions[idx])
	}
	return disjunction
}

func (binary Binary) Accept(visitor ExpressionVisitor) {
	visitor.Binary(binary)
}