package z3

// Represents an AST of the Bool sort.
type Bool struct {
	value bool
}

func NewBool(value bool) *Bool {
	return &Bool{
		value: value,
	}
}

func (operand *Bool) Not() *Bool {
	return NewBool(!operand.value)
}

func (lhs *Bool) And(rhs *Bool) *Bool {
	return NewBool(lhs.value && rhs.value)
}

func (lhs *Bool) Or(rhs *Bool) *Bool {
	return NewBool(lhs.value || rhs.value)
}

func (lhs *Bool) Xor(rhs *Bool) *Bool {
	return NewBool(lhs.value != rhs.value)
}

func (lhs *Bool) Eq(rhs *Bool) *Bool {
	return NewBool(lhs.value == rhs.value)
}

func (boolean *Bool) AST(context *Context) *AST {
	return newAST(context.bool(boolean.value))
}