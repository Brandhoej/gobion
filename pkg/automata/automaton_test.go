package automata

import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/automata/language/statements"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func TestCompletion(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[any](symbols.NewSymbolsFactory())
	x := expressions.NewVariable(symbols.Insert("x"), expressions.IntegerSort)

	builder := NewAutomatonBuilder()
	initial := builder.AddInitial("Initial")
	builder.AddLoop(initial,
		WithGuard(
			NewGuard(constraints.NewLogicalConstraint(
				expressions.NewBinary(
					x, expressions.LessThanEqual, expressions.NewInteger(10),
				),
			)),
		),
		WithUpdate(
			NewUpdate(
				constraints.NewAssignmentConstraint(
					statements.NewAssignment(
						x, expressions.NewBinary(
							x, expressions.Addition, expressions.NewInteger(1),
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
