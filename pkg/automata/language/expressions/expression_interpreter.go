package expressions

import "github.com/Brandhoej/gobion/internal/z3"

type SymbolicInterpreter struct {
	context   *z3.Context
	variables Variables
}

func NewSymbolicInterpreter(context *z3.Context, variables Variables) SymbolicInterpreter {
	return SymbolicInterpreter{
		context:    context,
		variables: variables,
	}
}

func (interpreter SymbolicInterpreter) Interpret(expression Expression) *z3.AST {
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
	}
	panic("Unknown expression type")
}

func (interpreter SymbolicInterpreter) Variable(variable Variable) *z3.AST {
	if sort, exists := interpreter.variables.Lookup(variable.symbol); exists {
		var z3Sort *z3.Sort
		switch sort {
		case BooleanSort:
			z3Sort = interpreter.context.BooleanSort()
		case IntegerSort:
			z3Sort = interpreter.context.IntegerSort()
		}
	
		return interpreter.context.NewConstant(
			z3.WithInt(int(variable.Symbol())), z3Sort,
		)
	}
	panic("Unknown variable")
}

func (interpreter SymbolicInterpreter) Binary(binary Binary) *z3.AST {
	lhs := interpreter.Interpret(binary.lhs)
	rhs := interpreter.Interpret(binary.rhs)
	switch binary.operator {
	case Equal:
		return z3.Eq(lhs, rhs)
	case NotEqual:
		return z3.Not(z3.Eq(lhs, rhs))
	case LessThan:
		return z3.LT(lhs, rhs)
	case LessThanEqual:
		return z3.LE(lhs, rhs)
	case GreaterThan:
		return z3.GT(lhs, rhs)
	case GreaterThanEqual:
		return z3.GE(lhs, rhs)
	case LogicalAnd:
		return z3.And(lhs, rhs)
	case LogicalOr:
		return z3.Or(lhs, rhs)
	case Addition:
		return z3.Add(lhs, rhs)
	case Subtraction:
		return z3.Subtract(lhs, rhs)
	}
	panic("Unknown binary operator")
}

func (interpreter SymbolicInterpreter) Integer(integer Integer) *z3.AST {
	return interpreter.context.NewInt(integer.value, interpreter.context.IntegerSort())
}

func (interpreter SymbolicInterpreter) Boolean(boolean Boolean) *z3.AST {
	if boolean.value {
		return interpreter.context.NewTrue()
	}
	return interpreter.context.NewFalse()
}

func (interpreter SymbolicInterpreter) Unary(unary Unary) *z3.AST {
	operand := interpreter.Interpret(unary.operand)
	switch unary.operator {
	case LogicalNegation:
		return z3.Not(operand)
	}
	panic("Unknown unary operator")
}
