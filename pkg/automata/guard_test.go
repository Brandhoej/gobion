package automata

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/automata/language/expressions"
	"github.com/Brandhoej/gobion/pkg/symbols"
	"github.com/stretchr/testify/assert"
)

func Test_IsSatisfiable(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())

	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	variables := expressions.NewVariablesMap[*z3.Sort]()
	x, y := symbols.Insert("x"), symbols.Insert("y")
	variables.Declare(x, context.IntegerSort())
	variables.Declare(y, context.IntegerSort())

	guard := NewGuard(
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
	)

	valuations := expressions.NewValuationsMap[*z3.AST]()
	solver := NewInterpreter(context, variables)

	for i := 0; i < 1000; i++ {
		// Act
		xVal, yVal := rand.Intn(1000)-500, rand.Intn(1000)-500
		valuations.Assign(x, context.NewInt(xVal, context.IntegerSort()))
		valuations.Assign(y, context.NewInt(yVal, context.IntegerSort()))
		satisfiable := guard.IsSatisfied(valuations, solver)

		// Assert
		expected := ((xVal >= 2) || (yVal <= 1 || yVal >= 3))
		if satisfiable != expected {
			assert.Equal(t, expected, satisfiable, fmt.Sprintf("Counter example with [x=%v, y=%v]", xVal, yVal))
		}
	}
}
