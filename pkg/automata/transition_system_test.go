package automata

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

func Test_Outgoing(t *testing.T) {
	// Arrange
	context := z3.NewContext(z3.NewConfig())
	zero := context.NewInt(0, context.IntegerSort())
	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())

	symbolsMap := symbols.NewSymbolsMap[string](
		symbols.NewSymbolsFactory(),
	)
	x := symbolsMap.Insert("x")
	variables := NewVariablesMap(context)
	xVar := variables.Declare(x, context.IntegerSort())
	valuations := NewValuationsMap(context)
	valuations.Assign(x, zero, EQ)

	locations := graph.NewVertexMap[Location]()
	start := locations.Add(NewLocation("Initial", NewInvariant(context, NewTrue(context))))
	final := locations.Add(NewLocation("Final", NewInvariant(context, NewTrue(context))))
	edges := graph.NewEdgesMap[Edge]()
	edges.Add(
		NewEdge(
			start,
			NewGuard(context, newEquality(xVar, GE, zero)),
			NewUpdate(context, NewAssignment(x, EQ, two)),
			final,
		),
	)
	edges.Add(
		NewEdge(
			start,
			NewGuard(context, NewFalse(context)),
			NewUpdate(context, NewAssignment(x, EQ, one)),
			final,
		),
	)
	dg := graph.NewLabeledDirected(locations, edges)
	automaton := NewAutomaton(context, dg, start)
	transitions := NewTransitionSystem(variables, automaton)

	// Act
	outgoings := transitions.Outgoing(transitions.Initial(valuations))

	// Assert
	outgoings[0].valuations.All(func(symbol symbols.Symbol, value Value) bool {
		variable, _ := variables.Variable(symbol)
		ast := value.relation.ast(variable, value.ast)
		t.Log(ast.String())
		return true
	})
	t.Log(outgoings)
	t.FailNow()
}
