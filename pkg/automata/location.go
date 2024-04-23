package automata

import (
	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
)

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

func (location Location) IsEnabled(valuations expressions.Valuations[*z3.AST], solver *ConstraintSolver) bool {
	return location.invariant.IsSatisfiable(valuations, solver)
}
