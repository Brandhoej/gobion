package cfg

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"strconv"
	"strings"
)

type DOT[S any, E comparable] struct {
	BlockIDs map[int]int
	ConditionIDs map[int]int
	statementStringer func(S) string
	expressionStringer func(E) string
}

func NewDOT[S any, E comparable](
	statementStringer func(S) string,
	expressionStringer func(E) string,
) *DOT[S, E] {
	return &DOT[S, E]{
		BlockIDs: map[int]int{},
		ConditionIDs: map[int]int{},
		statementStringer: statementStringer,
		expressionStringer: expressionStringer,
	}
}

func (dot *DOT[S, E]) Graph(writer io.Writer, graph *Graph[S, E]) {
	io.WriteString(writer, "digraph G {\n")
	dot.Nodes(writer, graph)
	dot.Edges(writer, graph)
	io.WriteString(writer, "}\n")
}

func (dot *DOT[S, E]) Nodes(writer io.Writer, graph *Graph[S, E]) {
	io.WriteString(
		writer, DotNode("initial", map[string]string{
			"shape": "point",
		}),
	)

	io.WriteString(
		writer, DotNode("terminal", map[string]string{
			"shape": "point",
		}),
	)

	for _, block := range graph.Blocks() {
		dot.Block(writer, graph, block)

		if _, condition := graph.Block(block); graph.IsConstrained(condition) {
			dot.Condition(writer, graph, condition)
		}
	}
}

func (dot *DOT[S, E]) Edges(writer io.Writer, graph *Graph[S, E]) {
	drawnJumps := make(map[int]bool)

	for _, block := range graph.Blocks() {
		blockID, _ := DotIDFor("block_", block, dot.BlockIDs)
		fromID := blockID

		_, condition := graph.Block(block)
		if graph.IsConstrained(condition) {
			fromID, _ = DotIDFor("condition_", condition, dot.ConditionIDs)
			io.WriteString(writer, DotEdge(blockID, fromID, nil))
		}

		_, jumps := graph.Condition(condition)
		for _, jump := range jumps {
			if drawn := drawnJumps[jump]; !drawn {
				dot.Jump(writer, graph, fromID, jump)
				drawnJumps[jump] = true
			}
		}
	}

	entryID, _ := DotIDFor("block_", graph.entry, dot.BlockIDs)
	io.WriteString(writer, DotEdge("initial", entryID, nil))
}

func (dot *DOT[S, E]) Block(writer io.Writer, graph *Graph[S, E], block int) {
	if blockID, first := DotIDFor("block_", block, dot.BlockIDs); first {
		statements, _ := graph.Block(block)
		// Write the block node without considering if the ID was intiially created.
		// We assume that the iteration of blocks will only be over unique blocks.
		blockAttributs := DotBlockAttributes(statements)
		io.WriteString(writer, DotNode(blockID, blockAttributs))
	}
}

func (dot *DOT[S, E]) Condition(writer io.Writer, graph *Graph[S, E], condition int) {
	// From the the start point for the edge to draw jumps from.
	// only if the block has a constrained condition will it be from the comdition.
	// This ensure that sequential executions of blocks is somewhat compacted.
	if graph.IsConstrained(condition) {
		// Since we have a condition node then we immediately draw the transition to it.
		// io.WriteString(writer, DotEdge(blockID, conditionID, nil))

		if conditionID, first := DotIDFor("condition_", condition, dot.ConditionIDs); first {
			expression, _ := graph.Condition(condition)
			conditionAttributes := DotConditionAttributes(expression)
			io.WriteString(writer, DotNode(conditionID, conditionAttributes))
		}
	}
}

func (dot *DOT[S, E]) Jump(writer io.Writer, graph *Graph[S, E], from string, jump int) {
	expression, destination := graph.Jump(jump)
	to := "terminal"
	if destination != graph.Exit() {
		to, _ = DotIDFor("block_", destination, dot.BlockIDs)
	}

	jumpAttributes := DotJumpAttributes(expression)
	io.WriteString(writer, DotEdge(from, to, jumpAttributes))
}

func DotAttributes(attributes map[string]string) string {
	var builder strings.Builder
	for key, value := range attributes {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%s=%s", key, value))
	}
	return builder.String()
}

func DotNodeAST(node any) string {
	if node == nil {
		return "\"\""
	}
	var buffer bytes.Buffer
	format.Node(&buffer, token.NewFileSet(), node)
	str := buffer.String()
	str = strings.ReplaceAll(str, "\"", "\"")
	str = strings.ReplaceAll(str, "\t", "")
	return strconv.Quote(str)
}

func DotNode(id any, attributes map[string]string) string {
	return fmt.Sprintf("%s [%s]\n", id, DotAttributes(attributes))
}

func DotEdge(from, to any, attributes map[string]string) string {
	return fmt.Sprintf("%s -> %s [%s]\n", from, to, DotAttributes(attributes))
}

func DotBlockAttributes(node any) map[string]string {
	return map[string]string{
		"label": DotNodeAST(node),
		"shape": "rectangle",
	}
}

func DotConditionAttributes(node any) map[string]string {
	return map[string]string{
		"label": DotNodeAST(node),
		"shape": "diamond",
	}
}

func DotJumpAttributes(node any) map[string]string {
	return map[string]string{
		"label": DotNodeAST(node),
	}
}

func DotIDFor[T comparable](prefix string, element T, ids map[T]int) (string, bool) {
	if id, found := ids[element]; found {
		return fmt.Sprintf("%s%v", prefix, id), false
	}
	id := len(ids)
	ids[element] = id
	return fmt.Sprintf("%s%v", prefix, id), true
}