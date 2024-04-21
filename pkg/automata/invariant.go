package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Invariant struct {
	constraint constraints.Constraint
}

func NewInvariant(constraint constraints.Constraint) Invariant {
	return Invariant{
		constraint: constraint,
	}
}

func NewTrueInvariant() Invariant {
	return NewInvariant(constraints.NewTrue())
}

func NewFalseInvariant() Invariant {
	return NewInvariant(constraints.NewFalse())
}

func (invariant Invariant) IsSatisfiable(valuations expressions.Valuations, solver *ConstraintSolver) bool {
	return solver.Satisfies(valuations, invariant.constraint)
}

func (invariant Invariant) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := constraints.NewPrettyPrinter(&buffer, symbols)
	printer.Constraint(invariant.constraint)
	return buffer.String()
}
