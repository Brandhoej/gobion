package automata

import "github.com/Brandhoej/gobion/pkg/automata/language"

type Guard struct {
	condition language.Expression
}

func NewGuard(condition language.Expression) Guard {
	return Guard{
		condition: condition,
	}
}

// Finds the union of all variables and adds conjunctive terms
func (guard Guard) Conjunction(guards ...Guard) Guard {
	conditions := make([]language.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := language.Conjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

// Finds the intersection of all variables and adds a disjunction term to them.
func (guard Guard) Disjunction(guards ...Guard) Guard {
	conditions := make([]language.Expression, len(guards))
	for idx := range guards {
		conditions[idx] = guards[idx].condition
	}
	conjunction := language.Disjunction(guard.condition, conditions...)
	return NewGuard(conjunction)
}

func (guard Guard) Negation() Guard {
	negation := language.LogicalNegate(guard.condition)
	return NewGuard(negation)
}

func (guard Guard) IsSatisfiable(solver Solver) bool {
	return solver.HasSolutionFor(guard.condition)
}

func (guard Guard) String() string {
	printer := language.NewPrettyPrinter()
	printer.Expression(guard.condition)
	return printer.String()
}
