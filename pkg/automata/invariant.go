package automata

import (
	"bytes"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Invariant struct {
	condition expressions.Expression
}

func NewInvariant(condition expressions.Expression) Invariant {
	return Invariant{
		condition: condition,
	}
}

func NewTrueInvariant() Invariant {
	return NewInvariant(expressions.NewTrue())
}

func NewFalseInvariant() Invariant {
	return NewInvariant(expressions.NewFalse())
}

func (invariant Invariant) IsSatisfiable(valuations expressions.Valuations[*z3.AST], solver *Interpreter) bool {
	return solver.IsSatisfied(valuations, invariant.condition)
}

func (invariant Invariant) String(symbols symbols.Store[any]) string {
	var buffer bytes.Buffer
	printer := expressions.NewPrettyPrinter(&buffer, symbols)
	printer.Expression(invariant.condition)
	return buffer.String()
}
