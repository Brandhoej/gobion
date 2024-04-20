package constraints

import (
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
)

type ConstraintVisitor interface {
	ExpressionConstraint(constraint ExpressionConstraint)
	AssignmentConstraint(constraint AssignmentConstraint)
	BinaryConstraint(constraint BinaryConstraint)
	UnaryConstraint(constraint UnaryConstraint)
}

type Constraint interface {
	Accept(visitor ConstraintVisitor)
}

type ExpressionConstraint struct {
	expression expressions.Expression
}

func NewLogicalConstraint(expression expressions.Expression) ExpressionConstraint {
	return ExpressionConstraint{
		expression: expression,
	}
}

func NewTrue() ExpressionConstraint {
	return NewLogicalConstraint(
		expressions.NewTrue(),
	)
}

func NewFalse() ExpressionConstraint {
	return NewLogicalConstraint(
		expressions.NewFalse(),
	)
}

func (constraint ExpressionConstraint) Accept(visitor ConstraintVisitor) {
	visitor.ExpressionConstraint(constraint)
}

type AssignmentConstraint struct {
	assignment statements.Assignment
}

func NewAssignmentConstraint(assignment statements.Assignment) AssignmentConstraint {
	return AssignmentConstraint{
		assignment: assignment,
	}
}

func (constraint AssignmentConstraint) Accept(visitor ConstraintVisitor) {
	visitor.AssignmentConstraint(constraint)
}

type BinaryConstraintOperator uint16

const (
	LogicalAnd = BinaryConstraintOperator(0)
	LogicalOr  = BinaryConstraintOperator(1)
	LogicalImplication = BinaryConstraintOperator(2)
)

type BinaryConstraint struct {
	lhs, rhs Constraint
	operator BinaryConstraintOperator
}

func NewBinary(
	lhs Constraint, operator BinaryConstraintOperator, rhs Constraint,
) BinaryConstraint {
	return BinaryConstraint{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

func Conjunction(constraint Constraint, constraints ...Constraint) (conjunction Constraint) {
	conjunction = constraint
	if len(constraints) == 0 {
		return conjunction
	}

	for idx := range constraints {
		conjunction = NewBinary(conjunction, LogicalAnd, constraints[idx])
	}
	return conjunction
}

func Disjunction(constraint Constraint, constraints ...Constraint) (disjunction Constraint) {
	disjunction = constraint
	if len(constraints) == 0 {
		return disjunction
	}

	for idx := range constraints {
		disjunction = NewBinary(disjunction, LogicalOr, constraints[idx])
	}
	return disjunction
}

func Implication(premise, conclusion Constraint) Constraint {
	return BinaryConstraint{
		lhs: premise,
		operator: LogicalImplication,
		rhs: conclusion,
	}
}

func (constraint BinaryConstraint) Accept(visitor ConstraintVisitor) {
	visitor.BinaryConstraint(constraint)
}

type UnaryConstraintOperator uint16

const (
	LogicalNegation = UnaryConstraintOperator(0)
)

type UnaryConstraint struct {
	operator UnaryConstraintOperator
	operand  Constraint
}

func NewUnaryConstraint(operator UnaryConstraintOperator, operand Constraint) UnaryConstraint {
	return UnaryConstraint{
		operator: operator,
		operand:  operand,
	}
}

func LogicalNegate(constraint Constraint) UnaryConstraint {
	return UnaryConstraint{
		operator: LogicalNegation,
		operand:  constraint,
	}
}

func (constraint UnaryConstraint) Accept(visitor ConstraintVisitor) {
	visitor.UnaryConstraint(constraint)
}
