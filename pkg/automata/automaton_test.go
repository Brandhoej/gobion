package automata

import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func TestCompletion(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[any](symbols.NewSymbolsFactory())
	x := symbols.Insert("x")

	builder := NewAutomatonBuilder()
	initial := builder.AddInitial("Initial")
	builder.AddLoop(initial,
		WithGuard(
			NewGuard(constraints.NewLogicalConstraint(
				expressions.NewBinary(
					expressions.NewValuation(x), expressions.LessThanEqual, expressions.NewInteger(10),
				),
			)),
		),
		WithUpdate(
			NewUpdate(
				constraints.NewLogicalConstraint(
					expressions.NewBinary(
						expressions.NewVariable(x),
						expressions.Equal,
						expressions.NewBinary(
							expressions.NewValuation(x),
							expressions.Addition,
							expressions.NewInteger(1),
						),
					),
				),
			),
		),
	)
	automaton := builder.Build(symbols)

	// Act
	var buffer bytes.Buffer
	automaton.DOT(&buffer)

	// Assert
	t.Log(buffer.String())
	t.FailNow()
}
