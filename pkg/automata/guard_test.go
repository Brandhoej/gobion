package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_IsSatisfiable(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	variables := language.NewVariablesMap()
	x := variables.Declare(symbols.Insert("x"), language.IntegerSort)
	y := variables.Declare(symbols.Insert("y"), language.IntegerSort)

	guard := NewGuard(
		language.Conjunction(
			language.NewBinary(
				x,
				language.Equal,
				language.NewInteger(2),
			),
			language.Disjunction(
				language.NewBinary(
					y,
					language.LessThanEqual,
					language.NewInteger(1),
				),
				language.NewBinary(
					y,
					language.GreaterThanEqual,
					language.NewInteger(3),
				),
			),
		),
	)

	valuations := language.NewValuationsMap()
	valuations.Assign(
		x.Symbol(),
		language.NewInteger(2),
	) // x = 2
	valuations.Assign(
		y.Symbol(),
		language.NewInteger(3),
	) // y = 1

	context := z3.NewContext(z3.NewConfig())
	solver := newSolver(context.NewSolver(), variables, valuations)

	// Act
	satisfiable := guard.IsSatisfiable(solver)

	// Assert
	t.Log(satisfiable)
	t.FailNow()
}
