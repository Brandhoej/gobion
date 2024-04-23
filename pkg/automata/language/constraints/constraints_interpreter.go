package constraints

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
)

type SymbolicInterpreter struct {
	context *z3.Context
	solver *z3.Solver
	expressions expressions.SymbolicInterpreter
	pc *z3.AST
}

func NewSymbolicInterpreter(
	context *z3.Context,
	valuations expressions.Valuations[*z3.AST],
	variables expressions.Variables[*z3.Sort],
) *SymbolicInterpreter {
	return &SymbolicInterpreter{
		context: context,
		solver: context.NewSolver(),
		expressions: expressions.NewSymbolicInterpreter(context, variables, valuations),
		pc: context.NewTrue(),
	}
}

func (interpreter *SymbolicInterpreter) Interpret(constraint Constraint) *z3.AST {
	switch cast := any(constraint).(type) {
	case ExpressionConstraint:
		return interpreter.ExpressionConstraint(cast)
	case BinaryConstraint:
		return interpreter.BinaryConstraint(cast)
	case UnaryConstraint:
		return interpreter.UnaryConstraint(cast)
	}
	panic("Unknown constraint type")
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