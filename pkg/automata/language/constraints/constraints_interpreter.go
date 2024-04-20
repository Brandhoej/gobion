package constraints

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
)

type SymbolicInterpreter struct {
	context *z3.Context
	expressions expressions.SymbolicInterpreter
	statements statements.SymbolicInterpreter
}

func NewSymbolicInterpreter(context *z3.Context, valuations expressions.Valuations) SymbolicInterpreter {
	return SymbolicInterpreter{
		context: context,
		expressions: expressions.NewSymbolicInterpreter(context, valuations),
		statements: statements.NewSymbolicInterpreter(context, valuations),
	}
}

func (interpreter SymbolicInterpreter) Interpret(constraint Constraint) *z3.AST {
	switch cast := any(constraint).(type) {
	case AssignmentConstraint:
		return interpreter.AssignmentConstraint(cast)
	case ExpressionConstraint:
		return interpreter.ExpressionConstraint(cast)
	case BinaryConstraint:
		return interpreter.BinaryConstraint(cast)
	case UnaryConstraint:
		return interpreter.UnaryConstraint(cast)
	}
	panic("Unknown constraint type")
}

func (interpreter SymbolicInterpreter) AssignmentConstraint(constraint AssignmentConstraint) *z3.AST {
	return interpreter.Statement(constraint.assignment)
}

func (interpreter SymbolicInterpreter) ExpressionConstraint(constraint ExpressionConstraint) *z3.AST {
	return interpreter.expressions.Interpret(constraint.expression)
}

func (interpreter SymbolicInterpreter) BinaryConstraint(constraint BinaryConstraint) *z3.AST {
	lhs := interpreter.Interpret(constraint.lhs)
	rhs := interpreter.Interpret(constraint.rhs)
	switch constraint.operator {
	case LogicalAnd:
		return z3.And(lhs, rhs)
	case LogicalOr:
		return z3.Or(lhs, rhs)
	case LogicalImplication:
		return z3.Implies(lhs, rhs)
	}
	panic("Unknown binary constraint operator")
}

func (interpreter SymbolicInterpreter) UnaryConstraint(constraint UnaryConstraint) *z3.AST {
	operand := interpreter.Interpret(constraint.operand)
	switch constraint.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary constraint operator")
}

func (interpreter SymbolicInterpreter) Statement(statement statements.Statement) *z3.AST {
	switch cast := any(statement).(type) {
	case statements.Assignment:
		return interpreter.Assignment(cast)
	}
	panic("Unknown statement type")
}

func (interpreter SymbolicInterpreter) Assignment(assignment statements.Assignment) *z3.AST {
	lhs := interpreter.expressions.Interpret(assignment.LHS())
	rhs := interpreter.expressions.Interpret(assignment.RHS())
	return z3.Eq(lhs, rhs)
}
