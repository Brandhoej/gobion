package cfg

import "io"

type block[S any, E comparable] struct {
	statements []S
	condition  int
}

type condition[S any, E comparable] struct {
	expression E
	jumps      []int
}

type jump[S any, E comparable] struct {
	expression  E
	destination int
}

type Graph[S any, E comparable] struct {
	blocks []*block[S, E]
	conditions []*condition[S, E]
	jumps []*jump[S, E]
	entry, exit  int
}

func New[S any, E comparable](statements ...S) *Graph[S, E] {
	graph := &Graph[S, E]{
		blocks: make([]*block[S, E], 0, 2),
		conditions: make([]*condition[S, E], 0, 2),
		jumps: make([]*jump[S, E], 0),
	}
	graph.entry, _ = graph.NewBlock(statements...)
	graph.exit = -1
	return graph
}

func (graph *Graph[S, E]) NewBlock(statements ...S) (int, int) {
	block := &block[S, E]{ statements: statements }
	graph.blocks = append(graph.blocks, block)
	id := len(graph.blocks)-1
	condition := graph.JumpTo(id, graph.NewUnconstainedCondition())
	return id, condition
}

func (graph *Graph[S, E]) NewCondition(expression E, jumps ...int) int {
	condition := &condition[S, E]{
		expression: expression,
		jumps:      jumps,
	}
	graph.conditions = append(graph.conditions, condition)
	return len(graph.conditions)-1
}

func (graph *Graph[S, E]) NewUnconstainedCondition(jumps ...int) int {
	var zero E
	return graph.NewCondition(zero, jumps...)
}

func (graph *Graph[S, E]) NewConditionalJump(expression E, destination int) int {
	jump := &jump[S, E]{
		expression:  expression,
		destination: destination,
	}
	graph.jumps = append(graph.jumps, jump)
	return len(graph.jumps)-1
}

func (graph *Graph[S, E]) NewUnconditionalJump(destination int) int {
	var zero E
	return graph.NewConditionalJump(zero, destination)
}

func (graph *Graph[S, E]) Blocks() []int {
	ids := make([]int, 0, len(graph.blocks))
	for id := range graph.blocks {
		ids = append(ids, id)
	}
	return ids
}

func (graph *Graph[S, E]) Append(block int, statments ...S) {
	graph.blocks[block].statements = append(graph.blocks[block].statements, statments...)
}

func (graph *Graph[S, E]) Block(id int) ([]S, int) {
	block := graph.blocks[id]
	return block.statements, block.condition
}

func (graph *Graph[S, E]) Condition(id int) (E, []int) {
	condition := graph.conditions[id]
	return condition.expression, condition.jumps
}

func (graph *Graph[S, E]) Jump(id int) (E, int) {
	jump := graph.jumps[id]
	return jump.expression, jump.destination
}

func (graph *Graph[S, E]) IsConstrained(condition int) bool {
	return !graph.IsUnconstrained(condition)
}

func (graph *Graph[S, E]) IsUnconstrained(condition int) bool {
	var zero E
	expression, _ := graph.Condition(condition)
	return expression == zero
}

func (graph *Graph[S, E]) Entry() int {
	return graph.entry
}

func (graph *Graph[S, E]) Exit() int {
	return graph.exit
}

func (graph *Graph[S, E]) Sequence(source, destination int) {
	jump := graph.NewUnconditionalJump(destination)
	condition := graph.NewUnconstainedCondition(jump)
	graph.JumpTo(source, condition)
}

func (graph *Graph[S, E]) IfThenElse(condition int, consequence, alternative int) {
	graph.JumpFrom(condition, consequence)
	graph.JumpFrom(condition, alternative)
}

func (graph *Graph[S, E]) JumpFrom(condition int, jump int) {
	a := graph.conditions[condition]
	a.jumps = append(a.jumps, jump)
}

func (graph *Graph[S, E]) JumpTo(source int, condition int) int {
	block := graph.blocks[source]
	block.condition = condition
	return condition
}

func (graph *Graph[S, E]) DOT(writer io.Writer) {
	NewDOT(
		func(s S) string {
			return DotNodeAST(s)
		},
		func(e E) string {
			return DotNodeAST(e)
		},
	).Graph(writer, graph)
}