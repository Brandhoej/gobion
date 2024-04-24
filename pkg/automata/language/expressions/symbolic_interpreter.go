package expressions

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/state"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type SymbolicInterpreter struct {
	context    *z3.Context
	variables  state.Variables[*z3.Sort]
	valuations state.Valuations[*z3.AST]
	pc *z3.AST
}

func NewSymbolicInterpreter(
	context *z3.Context,
	variables state.Variables[*z3.Sort],
	valuations state.Valuations[*z3.AST],
) *SymbolicInterpreter {
	return &SymbolicInterpreter{
		context:    context,
		variables:  variables,
		valuations: valuations,
		pc: context.NewTrue(),
	}
}

func (interpreter *SymbolicInterpreter) solver() *z3.Solver {
	solver := interpreter.context.NewSolver()
	solver.Assert(interpreter.pc)
	interpreter.valuations.All(func(symbol symbols.Symbol, value *z3.AST) bool {
		constant := interpreter.symbolVariable(symbol)
		equality := z3.Eq(constant, value)
		solver.Assert(equality)
		return true
	})
	return solver
}

func (interpreter *SymbolicInterpreter) isTrue(expression *z3.AST) bool {
	return interpreter.solver().Proven(z3.Eq(expression, interpreter.solver().True()))
}

func (interpreter *SymbolicInterpreter) isFalse(expression *z3.AST) bool {
	return interpreter.solver().Proven(z3.Eq(expression, interpreter.solver().False()))
}

func (interpreter *SymbolicInterpreter) canBeTrue(expression *z3.AST) bool {
	return !interpreter.isFalse(expression)
}

func (interpreter *SymbolicInterpreter) Satisfies(expression Expression) bool {
	interpretation := interpreter.Interpret(expression)
	return interpreter.isTrue(interpretation)
}

func (interpreter *SymbolicInterpreter) Interpret(expression Expression) *z3.AST {
	switch cast := any(expression).(type) {
	case Variable:
		return interpreter.variable(cast)
	case Binary:
		return interpreter.binary(cast)
	case Integer:
		return interpreter.integer(cast)
	case Boolean:
		return interpreter.boolean(cast)
	case Unary:
		return interpreter.unary(cast)
	case IfThenElse:
		return interpreter.ifThenElse(cast)
	case Assignment:
		return interpreter.assignment(cast)
	}
	panic("Unknown expression type")
}

func (interpreter *SymbolicInterpreter) symbolVariable(symbol symbols.Symbol) *z3.AST {
	if sort, exists := interpreter.variables.Lookup(symbol); exists {
		return interpreter.context.NewConstant(z3.WithInt(int(symbol)), sort)
	}
	return nil
}

func (interpreter *SymbolicInterpreter) variable(variable Variable) *z3.AST {
	if valuation, exists := interpreter.valuations.Value(variable.symbol); exists {
		return valuation
	}
	if constant := interpreter.symbolVariable(variable.symbol); constant != nil {
		return constant
	}
	panic("Unknown variable")
}

func (interpreter *SymbolicInterpreter) assignment(assignment Assignment) *z3.AST {
	valuation := interpreter.Interpret(assignment.valuation)
	symbol := assignment.variable.symbol
	interpreter.valuations.Assign(symbol, valuation)
	return interpreter.solver().True()
}

func (interpreter *SymbolicInterpreter) leftToRight(left, right Expression, shortCircuit func(lhs*z3.AST) bool) (*z3.AST, *z3.AST) {
	lhs := interpreter.Interpret(left)
	if shortCircuit != nil && shortCircuit(lhs) {
		return lhs, nil
	}
	rhs := interpreter.Interpret(right)
	return lhs, rhs
}

func (interpreter *SymbolicInterpreter) rightToLeft(left, right Expression) (*z3.AST, *z3.AST) {
	rhs := interpreter.Interpret(right)
	lhs := interpreter.Interpret(left)
	return lhs, rhs
}

func (interpreter *SymbolicInterpreter) binary(binary Binary) *z3.AST {
	switch binary.operator {
	case Equal:
		lhs, rhs := interpreter.rightToLeft(binary.lhs, binary.rhs)
		equality := z3.Eq(lhs, rhs)
		interpreter.pc = z3.And(interpreter.pc, equality)
		equal := interpreter.canBeTrue(equality)
		return interpreter.context.NewBoolean(equal)
	case NotEqual:
		lhs, rhs := interpreter.rightToLeft(binary.lhs, binary.rhs)
		equal := interpreter.canBeTrue(z3.Eq(lhs, rhs))
		return interpreter.context.NewBoolean(!equal)
	case LessThan:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		lessThan := z3.LT(lhs, rhs)
		return lessThan
	case LessThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		lessThanEqual := z3.LE(lhs, rhs)
		return lessThanEqual
	case GreaterThan:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		greaterThan := z3.GT(lhs, rhs)
		return greaterThan
	case GreaterThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		greaterThanEqual := z3.GE(lhs, rhs)
		return greaterThanEqual
	case Addition:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		return z3.Add(lhs, rhs)
	case Subtraction:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		return z3.Subtract(lhs, rhs)
	case LogicalAnd:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isFalse)
		if rhs == nil {
			return lhs
		}
		return z3.And(lhs, rhs)
	case LogicalOr:
		pc := interpreter.pc
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isTrue)
		if rhs == nil {
			interpreter.pc = pc
			return lhs
		}
		interpreter.pc = pc
		return z3.Or(lhs, rhs)
	case Implication:
		// P → Q ≡ (¬P) ∨ Q
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isFalse)
		if rhs == nil {
			return z3.Not(lhs)
		}
		return rhs
	}	
	panic("Unknown binary operator")
}

func (interpreter *SymbolicInterpreter) integer(integer Integer) *z3.AST {
	return interpreter.context.NewInt(integer.value, interpreter.context.IntegerSort())
}

func (interpreter *SymbolicInterpreter) boolean(boolean Boolean) *z3.AST {
	if boolean.value {
		return interpreter.context.NewTrue()
	}
	return interpreter.context.NewFalse()
}

func (interpreter *SymbolicInterpreter) unary(unary Unary) *z3.AST {
	operand := interpreter.Interpret(unary.operand)
	switch unary.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary operator")
}

func (interpreter *SymbolicInterpreter) ifThenElse(ite IfThenElse) *z3.AST {
	condition := interpreter.Interpret(ite.condition)
	if interpreter.isTrue(condition) {
		return interpreter.Interpret(ite.consequence)
	}
	return interpreter.Interpret(ite.alternative)
}
