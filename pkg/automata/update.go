package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Update struct {
	constraint constraints.Constraint
}

func NewUpdate(constraint constraints.Constraint) Update {
	return Update{
		constraint: constraint,
	}
}

func NewEmptyUpdate() Update {
	return Update{
		constraint: constraints.NewTrue(),
	}
}

func (update Update) Apply(solver *ConstraintSolver) {
	solver.Assert(update.constraint)
}

func (update Update) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := constraints.NewPrettyPrinter(&buffer, symbols)
	printer.Constraint(update.constraint)
	return buffer.String()
}