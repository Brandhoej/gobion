package automata

import "github.com/Brandhoej/gobion/pkg/automata/language"

type Invariant struct {
	condition language.Expression
}

func NewInvariant(condition language.Expression) Invariant {
	return Invariant{
		condition: condition,
	}
}

func (invariant Invariant) IsSatisfiable(solver Solver) bool {
	return solver.HasSolutionFor(invariant.condition)
}
