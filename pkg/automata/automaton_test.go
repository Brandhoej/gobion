package automata

import (
	"bytes"
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_DOT(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	symbols := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)

	// one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())

	variables := NewVariablesMap(context)
	x := symbols.Insert("x")

	variables.Declare(x, context.IntegerSort())

	guard := NewGuard(context, NewTrue(context))
	update := NewUpdate(
		context,
		NewAssignment(x, EQ, two),
	)

	locations := graph.NewVertexMap[Location]()
	loc0 := locations.Add(NewLocation("Start", NewInvariant(context, NewTrue(context))))
	loc1 := locations.Add(NewLocation("Action", NewInvariant(context, NewTrue(context))))
	locErr := locations.Add(NewLocation("Error", NewInvariant(context, NewTrue(context))))

	edges := graph.NewEdgesMap[Edge]()
	edges.Add(NewEdge(loc0, guard, update, loc1))

	dg := graph.NewLabeledDirected(locations, edges)

	automaton := NewAutomaton(context, dg, loc0)
	automaton.Complete(variables, func(k graph.Key, g Guard) graph.Key {
		// Angelic completion for initial location.
		if k == loc0 {
			return k
		}

		// All other use demonic completion.
		return locErr
	})

	// Act
	var buffer bytes.Buffer
	automaton.DOT(&buffer)

	// Assert
	t.Log(buffer.String())
	t.FailNow()
}
