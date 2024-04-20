package statements

import "github.com/Brandhoej/gobion/pkg/automata/language/expressions"

type StatementVisitor interface {
	Block(block Block)
	Assignment(assignment Assignment)
	ExpressionStatement(expressionStatement ExpressionStatement)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type Block struct {
	statements []Statement
}

func NewBlock(statements ...Statement) Block {
	return Block{
		statements: statements,
	}
}

func (block Block) Statements() []Statement {
	return block.statements
}

func (block Block) Accept(visitor StatementVisitor) {
	visitor.Block(block)
}

type Assignment struct {
	lhs expressions.Variable
	rhs expressions.Expression
}

func (assignment Assignment) LHS() expressions.Variable {
	return assignment.lhs
}

func (assignment Assignment) RHS() expressions.Expression {
	return assignment.rhs
}

func NewAssignment(lhs expressions.Variable, rhs expressions.Expression) Assignment {
	return Assignment{
		lhs: lhs,
		rhs: rhs,
	}
}

func (assignment Assignment) Accept(visitor StatementVisitor) {
	visitor.Assignment(assignment)
}

type ExpressionStatement struct {
	expression expressions.Expression
}

func NewExpressionStatement(expression expressions.Expression) ExpressionStatement {
	return ExpressionStatement{
		expression: expression,
	}
}

func (assignment ExpressionStatement) Accept(visitor StatementVisitor) {
	visitor.ExpressionStatement(assignment)
}
