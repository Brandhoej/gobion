package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_IsSatisfiable(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x, y := symbols.Insert("x"), symbols.Insert("y")

	guard := NewGuard(
		constraints.NewLogicalConstraint(
			expressions.Conjunction(
				expressions.NewBinary(
					expressions.NewVariable(x, expressions.IntegerSort),
					expressions.Equal,
					expressions.NewInteger(2),
				),
				expressions.Disjunction(
					expressions.NewBinary(
						expressions.NewVariable(y, expressions.IntegerSort),
						expressions.LessThanEqual,
						expressions.NewInteger(1),
					),
					expressions.NewBinary(
						expressions.NewVariable(y, expressions.IntegerSort),
						expressions.GreaterThanEqual,
						expressions.NewInteger(3),
					),
				),
			),
		),
	)

	context := z3.NewContext(z3.NewConfig())
	valuations := expressions.NewValuationsMap()
	variables := expressions.NewVariablesMap()
	solver := NewConstraintSolver(context.NewSolver(), valuations, variables)

	// Act
	satisfiable := guard.IsSatisfiable(solver)

	// Assert
	t.Log(satisfiable)
	t.FailNow()
}
