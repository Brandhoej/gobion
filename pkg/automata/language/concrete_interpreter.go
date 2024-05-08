package language

type ConcreteInterpreter struct {
	variables Variables
	valuations Valuations
}

func NewConcreteInterpreter(
	variables Variables,
	valuations Valuations,
) ConcreteInterpreter {
	return ConcreteInterpreter{
		variables: variables,
		valuations: valuations,
	}
}

func (interpreter ConcreteInterpreter) isTrue(expression Expression) bool {
	if boolean, ok := expression.(Boolean); ok {
		return boolean.Value()
	}
	panic("is true cannot be called on non-boolean expression")
}

func (interpreter ConcreteInterpreter) isFalse(expression Expression) bool {
	if boolean, ok := expression.(Boolean); ok {
		return !boolean.Value()
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
	}
	panic("Unknown expression type")
}

func (interpreter ConcreteInterpreter) variable(variable Variable) Expression {
	if valuation, exists := interpreter.valuations.Value(variable.Symbol()); exists {
		return valuation
	}
	panic("Varaible does not a valuation")
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
	switch binary.Operator() {
	case Equal:
		lhs, rhs := interpreter.rightToLeft(binary.LHS(), binary.RHS())

		// Equality between integers
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() == l.Value())
		}

		// Equality between booleans
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() == l.Value())
		}

		panic("Cannot perform EQ")
	case NotEqual:
		lhs, rhs := interpreter.rightToLeft(binary.LHS(), binary.RHS())

		// Equality between integers
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() != l.Value())
		}

		// Equality between booleans
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() != l.Value())
		}

		panic("Cannot perform NEQ")
	case LessThan:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() < l.Value())
		}

		panic("Cannot perform LT")
	case LessThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() <= l.Value())
		}

		panic("Cannot perform LE")
	case GreaterThan:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() > l.Value())
		}

		panic("Cannot perform GT")
	case GreaterThanEqual:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() >= l.Value())
		}

		panic("Cannot perform GE")
	case Addition:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewInteger(r.Value() + l.Value())
		}

		panic("Cannot perform ADD")
	case Subtraction:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), nil)
		if r, okR, l, okL := CastBinary[Integer, Integer](lhs, rhs); okR && okL {
			return NewInteger(r.Value() - l.Value())
		}

		panic("Cannot perform SUB")
	case LogicalAnd:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), interpreter.isFalse)
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() && l.Value())
		}

		panic("Cannot perform AND")
	case LogicalOr:
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), interpreter.isTrue)
		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(r.Value() || l.Value())
		}

		panic("Cannot perform OR")
	case Implication:
		// P → Q ≡ (¬P) ∨ Q
		lhs, rhs := interpreter.leftToRight(binary.LHS(), binary.RHS(), interpreter.isFalse)
		if rhs == nil {
			if l, ok := lhs.(Boolean); ok {
				return NewBoolean(!l.Value())
			}
		}

		if r, okR, l, okL := CastBinary[Boolean, Boolean](lhs, rhs); okR && okL {
			return NewBoolean(!l.Value() || r.Value())
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
	operand := interpreter.Interpret(unary.Operand())
	switch unary.Operator() {
	case LogicalNegation:
		if boolean, ok := operand.(Boolean); ok {
			return NewBoolean(!boolean.Value())
		}
	}
	panic("Unknown unary operator")
}

func (interpreter ConcreteInterpreter) ifThenElse(ite IfThenElse) Expression {
	condition := interpreter.Interpret(ite.Condition())
	if interpreter.isTrue(condition) {
		return interpreter.Interpret(ite.Consequence())
	}
	return interpreter.Interpret(ite.Alternative())
}
