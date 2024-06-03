package z3

import "math"

type Int struct {
	value int
}

func NewInt(value int) *Int {
	return &Int{
		value: value,
	}
}

func (operand *Int) Minus() *Int {
	return NewInt(-operand.value)
}

func (lhs *Int) Add(rhs *Int) *Int {
	return NewInt(lhs.value + rhs.value)
}

func (lhs *Int) Multiply(rhs *Int) *Int {
	return NewInt(lhs.value * rhs.value)
}

func (lhs *Int) Subtract(rhs *Int) *Int {
	return NewInt(lhs.value - rhs.value)
}

func (lhs *Int) Divide(rhs *Int) *Int {
	return NewInt(lhs.value / rhs.value)
}

func (lhs *Int) Modulus(rhs *Int) *Int {
	return NewInt(lhs.value % rhs.value)
}

func (lhs *Int) Remaninder(rhs *Int) *Int {
	remainder := lhs.value % rhs.value
	if remainder < 0 && rhs.value > 0 || remainder > 0 && rhs.value < 0 {
		remainder += rhs.value
	}
	return NewInt(remainder)
}

func (lhs *Int) Power(rhs *Int) *Int {
	return NewInt(int(math.Pow(float64(lhs.value), float64(rhs.value))))
}

func (lhs *Int) LT(rhs *Int) *Bool {
	return NewBool(lhs.value < rhs.value)
}

func (lhs *Int) LE(rhs *Int) *Bool {
	return NewBool(lhs.value <= rhs.value)
}

func (lhs *Int) GT(rhs *Int) *Bool {
	return NewBool(lhs.value > rhs.value)
}

func (lhs *Int) GE(rhs *Int) *Bool {
	return NewBool(lhs.value >= rhs.value)
}

func (lhs *Int) Divides(rhs *Int) *Bool {	
	return NewBool(lhs.value % rhs.value == 0)
}

func (integer *Int) AST(context *Context) *AST {
	return newAST(context.int(integer.value))
}