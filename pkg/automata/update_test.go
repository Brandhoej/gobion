package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_Apply(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())
	
	variables := NewVariablesMap(context)
	x := symbols.Insert("x")
	y := symbols.Insert("y")
	xVar := variables.Declare(x, context.IntegerSort())
	yVar := variables.Declare(y, context.IntegerSort())

	update := NewUpdate(context, NewAssignment(x, LE, z3.Add(two, yVar)))

	before := NewValuationsMap(context)
	before.Assign(x, one, EQ)
	before.Assign(y, two, EQ)
	
	// Act
	after := update.Apply(variables, before, true)

	// Assert
	a, _ := after.Value(x)
	t.Log(a.relation.ast(xVar, a.ast))
	t.FailNow()
}
