package language

type StatementVisitor interface {
	Assignment(assignment Assignment)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type Assignment struct {
	variable Variable
	valuation Expression
}

func NewAssignment(variable Variable, valuation Expression) Assignment {
	return Assignment{
		variable: variable,
		valuation: valuation,
	}
}

func (assignment Assignment) Accept(visitor StatementVisitor) {
	visitor.Assignment(assignment)
}