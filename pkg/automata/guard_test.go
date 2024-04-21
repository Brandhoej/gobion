package automata

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/constraints"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_IsSatisfiable(t *testing.T) {
	// Arrange
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	variables := expressions.NewVariablesMap()
	x, y := symbols.Insert("x"), symbols.Insert("y")
	variables.Declare(x, expressions.IntegerSort)
	variables.Declare(y, expressions.IntegerSort)

	guard := NewGuard(
		constraints.NewLogicalConstraint(
			expressions.Disjunction(
				expressions.NewBinary(
					expressions.NewVariable(x),
					expressions.GreaterThanEqual,
					expressions.NewInteger(2),
				),
				expressions.Disjunction(
					expressions.NewBinary(
						expressions.NewVariable(y),
						expressions.LessThanEqual,
						expressions.NewInteger(1),
					),
					expressions.NewBinary(
						expressions.NewVariable(y),
						expressions.GreaterThanEqual,
						expressions.NewInteger(3),
					),
				),
			),
		),
	)

	context := z3.NewContext(z3.NewConfig())
	valuations := expressions.NewValuationsMap()
	solver := NewConstraintSolver(context.NewSolver(), variables)

	for i := 0; i < 1000; i++ {
		// Act
		xVal, yVal := rand.Intn(1000)-500, rand.Intn(1000)-500
		valuations.Assign(x, expressions.NewInteger(xVal))
		valuations.Assign(y, expressions.NewInteger(yVal))
		satisfiable := guard.IsSatisfiable(valuations, solver)

		// Assert
		expected := ((xVal >= 2) || (yVal <= 1 || yVal >= 3))
		if satisfiable != expected {
			assert.Equal(t, expected, satisfiable, fmt.Sprintf("Counter example with [x=%v, y=%v]", xVal, yVal))
		}
	}
}