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
	valuations expressions.Valuations
}

func NewSymbolicInterpreter(
	context *z3.Context,
	valuations expressions.Valuations,
	variables expressions.Variables,
) *SymbolicInterpreter {
	return &SymbolicInterpreter{
		context: context,
		expressions: expressions.NewSymbolicInterpreter(context, variables),
		statements: statements.NewSymbolicInterpreter(context, valuations, variables),
		valuations: valuations,
	}
}

func (interpreter *SymbolicInterpreter) Interpret(constraint Constraint) *z3.AST {
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

func (interpreter *SymbolicInterpreter) AssignmentConstraint(constraint AssignmentConstraint) *z3.AST {
	return interpreter.Assignment(constraint.assignment)
}

func (interpreter *SymbolicInterpreter) Assignment(assignment statements.Assignment) *z3.AST {
	// TODO: This is actually incorrect when considering implications.
	// E.g., if (x >= 3) x=3 else x=1.
	// Symbolic execution should be performed on the constraint
	// with path merging to yield a valuation of "ITE"-form if necessary.
	valuation := interpreter.expressions.Interpret(assignment.RHS())
	interpreter.valuations.Assign(assignment.LHS().Symbol(), assignment.RHS())
	return valuation
}

func (interpreter *SymbolicInterpreter) ExpressionConstraint(constraint ExpressionConstraint) *z3.AST {
	return interpreter.expressions.Interpret(constraint.expression)
}

func (interpreter *SymbolicInterpreter) BinaryConstraint(constraint BinaryConstraint) *z3.AST {
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

func (interpreter *SymbolicInterpreter) UnaryConstraint(constraint UnaryConstraint) *z3.AST {
	operand := interpreter.Interpret(constraint.operand)
	switch constraint.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary constraint operator")
}