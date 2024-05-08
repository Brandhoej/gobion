package language

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type SymbolicInterpreter struct {
	context    *z3.Context
	variables  Variables
	valuations Valuations
	pc         *z3.AST
	z3solver   *z3.Solver
}

func NewSymbolicInterpreter(
	context *z3.Context,
	variables Variables,
	valuations Valuations,
) *SymbolicInterpreter {
	return &SymbolicInterpreter{
		context:    context,
		variables:  variables,
		valuations: valuations,
		pc:         context.NewTrue(),
		z3solver:   nil,
	}
}

func (interpreter *SymbolicInterpreter) solver() *z3.Solver {
	// Upon a change in valuations the backing z3 solver will be discarded.
	// However, if the solver exists then we can gurantee that it is valid.
	if interpreter.z3solver != nil {
		return interpreter.z3solver
	}

	interpreter.z3solver = interpreter.context.NewSolver()

	// An arbritary capacity is chosen. It will grow to fit the size of the valuation set.
	assertions := make([]*z3.AST, 1, 16)
	assertions[0] = interpreter.pc
	interpreter.valuations.All(func(symbol symbols.Symbol, value Expression) bool {
		constant := interpreter.symbolVariable(symbol)
		valuation := interpreter.Expression(value)
		equality := z3.Eq(constant, valuation)
		assertions = append(assertions, equality)
		return true
	})

	// Be deffering the assertion we can for valuation sets with many assignments
	// ommit many library calls to just a single one.
	if len(assertions) == 1 {
		interpreter.z3solver.Assert(assertions[0])
	} else if len(assertions) > 1 {
		interpreter.z3solver.Assert(z3.And(assertions[0], assertions[1:]...))
	}

	return interpreter.z3solver
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

func (interpreter *SymbolicInterpreter) sort(sort Sort) *z3.Sort {
	if sort == BooleanSort {
		return interpreter.context.BooleanSort()
	}
	if sort == IntegerSort {
		return interpreter.context.IntegerSort()
	}
	panic("Unknown sort")
}

func (interpreter *SymbolicInterpreter) Statement(statement Statement) {
	switch cast := any(statement).(type) {
	case Assignment:
		interpreter.Assignment(cast)
	default:
		panic("Unknown statement type")
	}
}

func (interpreter *SymbolicInterpreter) Assignment(assignment Assignment) {
	symbol := assignment.variable.Symbol()
	interpreter.valuations.Assign(symbol, assignment.valuation)
	interpreter.z3solver = nil
}

func (interpreter *SymbolicInterpreter) Expression(expression Expression) *z3.AST {
	switch cast := any(expression).(type) {
	case Variable:
		return interpreter.Variable(cast)
	case Binary:
		return interpreter.Binary(cast)
	case Integer:
		return interpreter.Integer(cast)
	case Boolean:
		return interpreter.Boolean(cast)
	case Unary:
		return interpreter.Unary(cast)
	case IfThenElse:
		return interpreter.IfThenElse(cast)
	case BlockExpression:
		return interpreter.BlockExpression(cast)
	}
	panic("Unknown expression type")
}

func (interpreter *SymbolicInterpreter) Satisfies(expression Expression) bool {
	interpretation := interpreter.Expression(expression)
	return interpreter.canBeTrue(interpretation)
}

func (interpreter *SymbolicInterpreter) symbolVariable(symbol symbols.Symbol) *z3.AST {
	if sort, exists := interpreter.variables.Lookup(symbol); exists {
		return interpreter.context.NewConstant(z3.WithInt(int(symbol)), interpreter.sort(sort))
	}
	return nil
}

func (interpreter *SymbolicInterpreter) Variable(variable Variable) *z3.AST {
	if value, exists := interpreter.valuations.Value(variable.symbol); exists {
		valuation := interpreter.Expression(value)
		return valuation
	}
	if constant := interpreter.symbolVariable(variable.symbol); constant != nil {
		return constant
	}
	panic("Unknown variable")
}

func (interpreter *SymbolicInterpreter) leftToRight(left, right Expression, shortCircuit func(lhs *z3.AST) bool) (*z3.AST, *z3.AST) {
	lhs := interpreter.Expression(left)
	if shortCircuit != nil && shortCircuit(lhs) {
		return lhs, nil
	}
	rhs := interpreter.Expression(right)
	return lhs, rhs
}

func (interpreter *SymbolicInterpreter) rightToLeft(left, right Expression) (*z3.AST, *z3.AST) {
	rhs := interpreter.Expression(right)
	lhs := interpreter.Expression(left)
	return lhs, rhs
}

func (interpreter *SymbolicInterpreter) Binary(binary Binary) *z3.AST {
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

func (interpreter *SymbolicInterpreter) Integer(integer Integer) *z3.AST {
	return interpreter.context.NewInt(integer.value, interpreter.context.IntegerSort())
}

func (interpreter *SymbolicInterpreter) Boolean(boolean Boolean) *z3.AST {
	if boolean.value {
		return interpreter.context.NewTrue()
	}
	return interpreter.context.NewFalse()
}

func (interpreter *SymbolicInterpreter) Unary(unary Unary) *z3.AST {
	operand := interpreter.Expression(unary.operand)
	switch unary.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary operator")
}

func (interpreter *SymbolicInterpreter) IfThenElse(ite IfThenElse) *z3.AST {
	condition := interpreter.Expression(ite.condition)
	if interpreter.isTrue(condition) {
		return interpreter.Expression(ite.consequence)
	}
	return interpreter.Expression(ite.alternative)
}

func (interpreter *SymbolicInterpreter) BlockExpression(blockExpression BlockExpression) *z3.AST {
	for idx := range blockExpression.statements {
		interpreter.Statement(blockExpression.statements[idx])
	}
	return interpreter.Expression(blockExpression.expression)
}
