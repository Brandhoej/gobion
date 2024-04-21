package automata

import "github.com/Brandhoej/gobion/pkg/automata/language/expressions"

type Location struct {
	name      string
	invariant Invariant
}

func NewLocation(name string, invariant Invariant) Location {
	return Location{
		name:      name,
		invariant: invariant,
	}
}

func (location Location) IsEnabled(valuations expressions.Valuations, solver *ConstraintSolver) bool {
	return location.invariant.IsSatisfiable(valuations, solver)
}
