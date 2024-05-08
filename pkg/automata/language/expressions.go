package language

import "github.com/Brandhoej/gobion/pkg/symbols"

type ExpressionVisitor interface {
	Variable(variable Variable)
	Binary(binary Binary)
	Integer(integer Integer)
	Boolean(boolean Boolean)
	Unary(unary Unary)
	IfThenElse(ite IfThenElse)
	BlockExpression(block BlockExpression)
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

type Variable struct {
	symbol symbols.Symbol
}

func NewVariable(symbol symbols.Symbol) Variable {
	return Variable{
		symbol: symbol,
	}
}

func (variable Variable) Symbol() symbols.Symbol {
	return variable.symbol
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
	LogicalOr        = BinaryOperator(7)
	Addition         = BinaryOperator(8)
	Subtraction      = BinaryOperator(9)
	Implication      = BinaryOperator(10)
)

type Binary struct {
	lhs, rhs Expression
	operator BinaryOperator
}

func NewBinary(lhs Expression, operator BinaryOperator, rhs Expression) Binary {
	return Binary{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

func (binary Binary) LHS() Expression {
	return binary.lhs
}

func (binary Binary) Operator() BinaryOperator {
	return binary.operator
}

func (binary Binary) RHS() Expression {
	return binary.rhs
}

func CastBinary[L, R any](lhs, rhs Expression) (l L, okL bool, r R, okR bool) {
	l, okL = lhs.(L)
	r, okR = rhs.(R)
	return l, okL, r, okR
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

type UnaryOperator uint16

const (
	LogicalNegation = UnaryOperator(0)
)

type Unary struct {
	operator UnaryOperator
	operand  Expression
}

func NewUnary(operator UnaryOperator, operand Expression) Unary {
	return Unary{
		operator: operator,
		operand:  operand,
	}
}

func (unary Unary) Operator() UnaryOperator {
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

type IfThenElse struct {
	condition, consequence, alternative Expression
}

func NewIfThenElse(condition, consequence, alternative Expression) IfThenElse {
	return IfThenElse{
		condition:   condition,
		consequence: consequence,
		alternative: alternative,
	}
}

func (ite IfThenElse) Condition() Expression {
	return ite.condition
}

func (ite IfThenElse) Consequence() Expression {
	return ite.consequence
}

func (ite IfThenElse) Alternative() Expression {
	return ite.alternative
}

func (ite IfThenElse) Accept(visitor ExpressionVisitor) {
	visitor.IfThenElse(ite)
}

type BlockExpression struct {
	statements []Statement
	expression Expression
}

func NewBlockExpression(expression Expression, statements ...Statement) BlockExpression {
	return BlockExpression{
		statements: statements,
		expression: expression,
	}
}

func (block BlockExpression) Accept(visitor ExpressionVisitor) {
	visitor.BlockExpression(block)
}