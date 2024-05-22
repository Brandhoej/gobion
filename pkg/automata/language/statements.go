package language

type StatementVisitor interface {
	Assignment(assignment Assignment)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type Assignment struct {
	lhs Expression
	rhs Expression
}

func NewAssignment(variable Expression, valuation Expression) Assignment {
	return Assignment{
		lhs: variable,
		rhs: valuation,
	}
}

func (assignment Assignment) Accept(visitor StatementVisitor) {
	visitor.Assignment(assignment)
}
