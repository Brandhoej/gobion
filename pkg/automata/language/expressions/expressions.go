package expressions

import "github.com/Brandhoej/gobion/pkg/symbols"

type ExpressionVisitor interface {
	Variable(variable Variable)
	Binary(binary Binary)
	Integer(integer Integer)
	Boolean(boolean Boolean)
	Unary(unary Unary)
}

type Expression interface {
	Accept(visitor ExpressionVisitor)
}

type Boolean struct {
	value bool
}

func NewBoolean(value bool) Boolean {
	return Boolean{
		value: value,
	}
}

func (boolean Boolean) Value() bool {
	return boolean.value
}

func NewFalse() Boolean {
	return NewBoolean(false)
}

func NewTrue() Boolean {
	return NewBoolean(true)
}

func (boolean Boolean) Accept(visitor ExpressionVisitor) {
	visitor.Boolean(boolean)
}

type Integer struct {
	value int
}

func NewInteger(value int) Integer {
	return Integer{
		value: value,
	}
}

func (integer Integer) Value() int {
	return integer.value
}

func (integer Integer) Accept(visitor ExpressionVisitor) {
	visitor.Integer(integer)
}

type Sort uint16

const (
	BooleanSort = Sort(0)
	IntegerSort = Sort(1)
)

type Variable struct {
	symbol symbols.Symbol
	sort   Sort
}

func NewVariable(symbol symbols.Symbol, sort Sort) Variable {
	return Variable{
		symbol: symbol,
		sort:   sort,
	}
}

func (variable Variable) Sort() Sort {
	return variable.sort
}

func (variable Variable) Symbol() symbols.Symbol {
	return variable.symbol
}

func (variable Variable) Accept(visitor ExpressionVisitor) {
	visitor.Variable(variable)
}

type BinaryExpressionOperator uint16

const (
	Equal            = BinaryExpressionOperator(0)
	NotEqual         = BinaryExpressionOperator(1)
	LessThan         = BinaryExpressionOperator(2)
	LessThanEqual    = BinaryExpressionOperator(3)
	GreaterThan      = BinaryExpressionOperator(4)
	GreaterThanEqual = BinaryExpressionOperator(5)
	LogicalAnd       = BinaryExpressionOperator(6)
	LogicalOr        = BinaryExpressionOperator(7)
	Addition         = BinaryExpressionOperator(8)
	Subtraction      = BinaryExpressionOperator(9)
)

type Binary struct {
	lhs, rhs Expression
	operator BinaryExpressionOperator
}

func NewBinary(lhs Expression, operator BinaryExpressionOperator, rhs Expression) Binary {
	return Binary{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

func (binary Binary) LHS() Expression {
	return binary.lhs
}

func (binary Binary) Operator() BinaryExpressionOperator {
	return binary.operator
}

func (binary Binary) RHS() Expression {
	return binary.rhs
}

func Conjunction(expression Expression, expressions ...Expression) (conjunction Expression) {
	conjunction = expression
	if len(expressions) == 0 {
		return conjunction
	}

	for idx := range expressions {
		conjunction = NewBinary(conjunction, LogicalAnd, expressions[idx])
	}
	return conjunction
}

func Disjunction(expression Expression, expressions ...Expression) (disjunction Expression) {
	disjunction = expression
	if len(expressions) == 0 {
		return disjunction
	}

	for idx := range expressions {
		disjunction = NewBinary(disjunction, LogicalOr, expressions[idx])
	}
	return disjunction
}

type UnaryExpressionOperator uint16

const (
	LogicalNegation = UnaryExpressionOperator(0)
)

type Unary struct {
	operator UnaryExpressionOperator
	operand  Expression
}

func NewUnary(operator UnaryExpressionOperator, operand Expression) Unary {
	return Unary{
		operator: operator,
		operand:  operand,
	}
}

func (unary Unary) Operator() UnaryExpressionOperator {
	return unary.operator
}

func (unary Unary) Operand() Expression {
	return unary.operand
}

func (unary Unary) Accept(visitor ExpressionVisitor) {
	visitor.Unary(unary)
}

func LogicalNegate(expression Expression) Unary {
	return NewUnary(LogicalNegation, expression)
}

func (binary Binary) Accept(visitor ExpressionVisitor) {
	visitor.Binary(binary)
}
