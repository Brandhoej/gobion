package statements

import "github.com/Brandhoej/gobion/pkg/automata/language/expressions"

type StatementVisitor interface {
	Assignment(assignment Assignment)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type Assignment struct {
	variable expressions.Variable
	valuation expressions.Expression
}

func NewAssignment(variable expressions.Variable, valuation expressions.Expression) Assignment {
	return Assignment{
		variable: variable,
		valuation: valuation,
	}
}

func (assignment Assignment) Accept(visitor StatementVisitor) {
	visitor.Assignment(assignment)
}