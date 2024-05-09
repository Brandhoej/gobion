package scfg

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
)

type DOT[S any, E comparable] struct {
	blockIDs           map[int]int
	conditionIDs       map[int]int
	statementStringer  func(S) string
	expressionStringer func(E) string
	FlowDOT            *cfg.DOT[S, E]
}

func NewDOT[S any, E comparable](
	statementStringer func(S) string,
	expressionStringer func(E) string,
) *DOT[S, E] {
	return &DOT[S, E]{
		statementStringer:  statementStringer,
		expressionStringer: expressionStringer,
		FlowDOT:            cfg.NewDOT(statementStringer, expressionStringer),
	}
}

func (dot *DOT[S, E]) Graph(writer io.Writer, graph *Graph[S, E]) {
	io.WriteString(writer, "digraph G {\n")
	dot.Nodes(writer, graph.Global(), graph)
	dot.FlowDOT.Edges(writer, graph.CFG())
	io.WriteString(writer, "}\n")
}

func (dot *DOT[S, E]) Nodes(writer io.Writer, scope int, graph *Graph[S, E]) {
	io.WriteString(
		writer, cfg.DotNode("initial", map[string]string{
			"shape": "point",
		}),
	)

	io.WriteString(
		writer, cfg.DotNode("terminal", map[string]string{
			"shape": "point",
		}),
	)

	var recursive func(scope int)
	recursive = func(scope int) {
		io.WriteString(writer, fmt.Sprintf("subgraph cluster_%v {\n", scope))

		for _, block := range graph.Blocks(scope) {
			flowGraph := graph.CFG()
			dot.FlowDOT.Block(writer, flowGraph, block)

			_, condition := flowGraph.Block(block)
			if flowGraph.IsConstrained(condition) {
				dot.FlowDOT.Condition(writer, flowGraph, condition)
			}
		}

		for _, child := range graph.Children(scope) {
			recursive(child)
		}

		io.WriteString(writer, "}\n")
	}
	recursive(scope)
}
