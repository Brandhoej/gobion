package language

import (
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/Brandhoej/gobion/pkg/zones"
)

type StatementVisitor interface {
	Assignment(assignment Assignment)
	ClockConstraint(constraint ClockConstraint)
	ClockAssignment(assignment ClockAssignment)
	ClockShift(shift ClockShift)
	ClockReset(reset ClockReset)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type ClockReset struct {
	clock symbols.Symbol
	limit int
}

func NewClockReset(clock symbols.Symbol, limit int) ClockReset {
	return ClockReset{
		clock: clock,
		limit: limit,
	}
}

func (free ClockReset) Accept(visitor StatementVisitor) {
	visitor.ClockReset(free)
}

type ClockShift struct {
	clock symbols.Symbol
	limit int
}

func NewClockShift(clock symbols.Symbol, limit int) ClockShift {
	return ClockShift{
		clock: clock,
		limit: limit,
	}
}

func (shift ClockShift) Accept(visitor StatementVisitor) {
	visitor.ClockShift(shift)
}

type ClockAssignment struct {
	lhs, rhs symbols.Symbol
}

func NewClockAssignment(lhs, rhs symbols.Symbol) ClockAssignment {
	return ClockAssignment{
		lhs: lhs,
		rhs: rhs,
	}
}

func (assignment ClockAssignment) Accept(visitor StatementVisitor) {
	visitor.ClockAssignment(assignment)
}

type ClockConstraint struct {
	lhs, rhs symbols.Symbol
	relation zones.Relation
}

func NewClockConstraint(lhs, rhs symbols.Symbol, relation zones.Relation) ClockConstraint {
	return ClockConstraint{
		lhs:      lhs,
		rhs:      rhs,
		relation: relation,
	}
}

func (constraint ClockConstraint) Accept(visitor StatementVisitor) {
	visitor.ClockConstraint(constraint)
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
