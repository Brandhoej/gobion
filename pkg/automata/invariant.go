package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Invariant struct {
	condition language.Expression
}

func NewInvariant(condition language.Expression) Invariant {
	return Invariant{
		condition: condition,
	}
}

func NewTrueInvariant() Invariant {
	return NewInvariant(language.NewTrue())
}

func NewFalseInvariant() Invariant {
	return NewInvariant(language.NewFalse())
}

func (invariant Invariant) IsSatisfiable(valuations language.Valuations, solver *Interpreter) bool {
	return solver.IsSatisfied(valuations, invariant.condition)
}

func (invariant Invariant) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := language.NewPrettyPrinter(&buffer, symbols)
	printer.Expression(invariant.condition)
	return buffer.String()
}
