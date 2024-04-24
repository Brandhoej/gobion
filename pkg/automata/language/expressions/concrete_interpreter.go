package expressions

import "github.com/Brandhoej/gobion/pkg/automata/language/state"

type ConcreteInterpreter struct {
	variables state.Variables[Sort]
	valuations state.Valuations[Expression]
}

func NewConcreteInterpreter(
	variables state.Variables[Sort],
	valuations state.Valuations[Expression],
) ConcreteInterpreter {
	return ConcreteInterpreter{
		variables: variables,
		valuations: valuations,
	}
}

func (interpreter ConcreteInterpreter) isTrue(expression Expression) bool {
	if boolean, ok := expression.(Boolean); ok {
		return boolean.value
	}
	panic("is true cannot be called on non-boolean expression")
}

func (interpreter ConcreteInterpreter) isFalse(expression Expression) bool {
	if boolean, ok := expression.(Boolean); ok {
		return !boolean.value
	}
	panic("is true cannot be called on non-boolean expression")
}

func (interpreter ConcreteInterpreter) Satisfies(expression Expression) bool {
	return interpreter.isTrue(interpreter.Interpret(expression))
}

func (interpreter ConcreteInterpreter) Interpret(expression Expression) Expression {
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

func (interpreter ConcreteInterpreter) variable(variable Variable) Expression {
	if valuation, exists := interpreter.valuations.Value(variable.symbol); exists {
		return valuation
	}
	panic("Varaible does not a valuation")
}

func (interpreter ConcreteInterpreter) assignment(assignment Assignment) Expression {
	valuation := interpreter.Interpret(assignment.valuation)
	symbol := assignment.variable.symbol
	interpreter.valuations.Assign(symbol, valuation)
	return NewBoolean(true)
}

func (interpreter ConcreteInterpreter) leftToRight(left, right Expression, shortCircuit func(Expression) bool) (Expression, Expression) {
	lhs := interpreter.Interpret(left)
	if shortCircuit != nil && shortCircuit(lhs) {
		return lhs, nil
	}
	rhs := interpreter.Interpret(right)
	return lhs, rhs
}

func (interpreter ConcreteInterpreter) rightToLeft(left, right Expression) (Expression, Expression) {
	rhs := interpreter.Interpret(right)
	lhs := interpreter.Interpret(left)
	return lhs, rhs
}

func (interpreter ConcreteInterpreter) binary(binary Binary) Expression {
	switch binary.operator {
	case Equal:
		lhs, rhs := interpreter.rightToLeft(binary.lhs, binary.rhs)

		// Equality between integers
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value == l.value)
		}

		// Equality between booleans
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.value == l.value)
		}

		panic("Cannot perform EQ")
	case NotEqual:
		lhs, rhs := interpreter.rightToLeft(binary.lhs, binary.rhs)

		// Equality between integers
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value != l.value)
		}

		// Equality between booleans
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.value != l.value)
		}

		panic("Cannot perform NEQ")
	case LessThan:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value < l.value)
		}

		panic("Cannot perform LT")
	case LessThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value <= l.value)
		}

		panic("Cannot perform LE")
	case GreaterThan:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value > l.value)
		}

		panic("Cannot perform GT")
	case GreaterThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.value >= l.value)
		}

		panic("Cannot perform GE")
	case Addition:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewInteger(r.value + l.value)
		}

		panic("Cannot perform ADD")
	case Subtraction:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewInteger(r.value - l.value)
		}

		panic("Cannot perform SUB")
	case LogicalAnd:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isFalse)
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.value && l.value)
		}

		panic("Cannot perform AND")
	case LogicalOr:
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isTrue)
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.value || l.value)
		}

		panic("Cannot perform OR")
	case Implication:
		// P → Q ≡ (¬P) ∨ Q
		lhs, rhs := interpreter.leftToRight(binary.lhs, binary.rhs, interpreter.isFalse)
		if rhs == nil {
			if l, ok := lhs.(Boolean); ok {
				return NewBoolean(!l.value)
			}
		}

		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(!l.value || r.value)
		}

		panic("Cannot perform OR")
	}	
	panic("Unknown binary operator")
}

func (interpreter ConcreteInterpreter) integer(integer Integer) Expression {
	return integer
}

func (interpreter ConcreteInterpreter) boolean(boolean Boolean) Expression {
	return boolean
}

func (interpreter ConcreteInterpreter) unary(unary Unary) Expression {
	operand := interpreter.Interpret(unary.operand)
	switch unary.operator {
	case LogicalNegation:
		if boolean, ok := operand.(Boolean); ok {
			return NewBoolean(!boolean.value)
		}
	}
	panic("Unknown unary operator")
}

func (interpreter ConcreteInterpreter) ifThenElse(ite IfThenElse) Expression {
	condition := interpreter.Interpret(ite.condition)
	if interpreter.isTrue(condition) {
		return interpreter.Interpret(ite.consequence)
	}
	return interpreter.Interpret(ite.alternative)
}
