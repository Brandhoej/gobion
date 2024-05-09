package language

import (
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/Brandhoej/gobion/pkg/zones"
)

type StatementVisitor interface {
	Assignment(assignment Assignment)
	ClockAssignment(assignment ClockAssignment)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type ClockAssignment struct {
	lhs, rhs symbols.Symbol
	relation zones.Relation
}

func NewClockAssignment(lhs, rhs symbols.Symbol, relation zones.Relation) ClockAssignment {
	return ClockAssignment{
		lhs:      lhs,
		rhs:      rhs,
		relation: relation,
	}
}

func (assignment ClockAssignment) Accept(visitor StatementVisitor) {
	visitor.ClockAssignment(assignment)
}

type Assignment struct {
	lhs Expression
	rhs Expression
}

func NewAssignment(variable Expression, valuation Expression) Assignment {
	return Assignment{
		lhs: variable,
		rhs: valuation,
	}
}

func (assignment Assignment) Accept(visitor StatementVisitor) {
	visitor.Assignment(assignment)
}
