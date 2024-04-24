package automata

import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func TestCompletion(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	symbols := symbols.NewSymbolsMap[any](symbols.NewSymbolsFactory())
	x := symbols.Insert("x")
	variables := expressions.NewVariablesMap[*z3.Sort]()
	variables.Declare(x, context.IntegerSort())

	builder := NewAutomatonBuilder()
	initial := builder.AddInitial("Initial")
	builder.AddLoop(initial,
		WithGuard(
			NewGuard(
				expressions.NewBinary(
					expressions.NewVariable(x), expressions.LessThanEqual, expressions.NewInteger(10),
				),
			),
		),
		WithUpdate(
			NewUpdate(
				expressions.NewAssignment(
					expressions.NewVariable(x),
					expressions.NewBinary(
						expressions.NewVariable(x),
						expressions.Addition,
						expressions.NewInteger(1),
					),
				),
			),
		),
	)
	err := builder.AddLocation("Error", WithInvariant(NewFalseInvariant()))
	automaton := builder.Build(symbols)
	automaton.Complete(NewInterpreter(context, variables), DirectedCompletion(err))

	// Act
	var buffer bytes.Buffer
	automaton.DOT(&buffer)

	// Assert
	t.Log(buffer.String())
	t.FailNow()
}
