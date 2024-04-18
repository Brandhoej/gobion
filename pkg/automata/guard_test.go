package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_IsSatisfiable(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())
	three := context.NewInt(2, context.IntegerSort())

	variables := NewVariablesMap(context)
	x := symbols.Insert("x")
	y := symbols.Insert("y")
	xVar := variables.Declare(x, context.IntegerSort())
	yVar := variables.Declare(y, context.IntegerSort())

	guard := NewGuard(
		context,
		NewConjunction(
			newEquality(xVar, EQ, two),
			NewDisjunction(
				newEquality(yVar, LE, one),
				newEquality(yVar, GE, three),
			),
		),
	)

	valuations := NewValuationsMap(context)
	valuations.Assign(x, two, LE) // x <= 2
	valuations.Assign(y, one, EQ) // y = 1

	// Act
	satisfiable := guard.IsSatisfiable(variables, valuations)

	// Assert
	t.Log(satisfiable)
	t.FailNow()
}
