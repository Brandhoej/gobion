package automata

import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func TestCompletion(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	symbols := symbols.NewSymbolsMap[any](symbols.NewSymbolsFactory())
	x := symbols.Insert("x")
	variables := language.NewVariablesMap()
	variables.Declare(x, language.IntegerSort)

	builder := NewAutomatonBuilder()
	initial := builder.AddInitial("Initial")
	builder.AddLoop(initial,
		WithGuard(
			NewGuard(
				language.NewBinary(
					language.NewVariable(x), language.LessThanEqual, language.NewInteger(10),
				),
			),
		),
		WithUpdate(
			NewUpdate(
				language.NewBlockExpression(
					language.NewTrue(),
					language.NewAssignment(
						language.NewVariable(x),
						language.NewBinary(
							language.NewVariable(x),
							language.Addition,
							language.NewInteger(1),
						),
					),
				),
			),
		),
	)
	err := builder.AddLocation("Error", WithInvariant(NewFalseInvariant()))
	automaton := builder.Build()
	automaton.Complete(NewInterpreter(context, variables), DirectedCompletion(err))

	// Act
	var buffer bytes.Buffer
	automaton.DOT(&buffer, symbols)

	// Assert
	t.Log(buffer.String())
	t.FailNow()
}
