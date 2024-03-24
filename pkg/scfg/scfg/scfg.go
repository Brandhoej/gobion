package scfg

import (
	"io"

	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
)

type scope[S any, E comparable] struct {
	parent   int
	children []int
}

type Graph[S any, E comparable] struct {
	flow   *cfg.Graph[S, E]
	scopes []*scope[S, E]
	blocks map[int]int
	global int
}

func New[S any, E comparable](flow *cfg.Graph[S, E]) *Graph[S, E] {
	global := &scope[S, E]{
		parent: 0,
		children: make([]int, 0),
	}
	return &Graph[S, E]{
		flow:   flow,
		scopes: []*scope[S, E]{global},
		blocks: map[int]int{},
		global: 0,
	}
}

func (graph *Graph[S, E]) CFG() *cfg.Graph[S, E] {
	return graph.flow
}

func (graph *Graph[S, E]) Global() int {
	return graph.global
}

func (graph *Graph[S, E]) Into(scope int, blocks ...int) {
	for _, block := range blocks {
		graph.blocks[block] = scope
	}
}

func (graph *Graph[S, E]) Blocks(scope int) (blocks []int) {
	for block, in := range graph.blocks {
		if in == scope {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (graph *Graph[S, E]) Parent(id int) (parent int, global bool) {
	if id == graph.global {
		return 0, true
	}
	scope := graph.scopes[id]
	return scope.parent, false
}

func (graph *Graph[S, E]) Children(id int) []int {
	scope := graph.scopes[id]
	return scope.children
}

func (graph *Graph[S, E]) ScopeWith(block int) (int, bool) {
	scope, exists := graph.blocks[block]
	return scope, exists
}

func (graph *Graph[S, E]) ZoomIn(parent int) int {
	child := &scope[S, E]{
		parent: parent,
		children: make([]int, 0),
	}
	graph.scopes = append(graph.scopes, child)
	childID := len(graph.scopes)-1
	scope := graph.scopes[parent]
	scope.children = append(scope.children, childID)
	return childID
}

func (graph *Graph[S, E]) ZoomOut(scope int) (int, bool) {
	return graph.Parent(scope)
}

func (graph *Graph[S, E]) Transitive(parent int) (blocks []int) {
	visited := make(map[int]bool)

	var helper func(scope int)
	helper = func(scope int) {
		if _, found := visited[scope]; found {
			return
		}

		blocks = append(blocks, graph.Blocks(scope)...)
		for _, child := range graph.Children(scope) {
			helper(child)
		}
	}

	helper(parent)

	return
}

func (graph *Graph[S, E]) DOT(writer io.Writer) {
	NewDOT(
		func(s S) string {
			return cfg.DotNodeAST(s)
		},
		func(e E) string {
			return cfg.DotNodeAST(e)
		},
	).Graph(writer, graph)
}