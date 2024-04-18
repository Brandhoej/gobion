package automata

import (
	"github.com/Brandhoej/gobion/pkg/algorithms"
	"github.com/Brandhoej/gobion/pkg/structures"
)

type SearchStrategy interface {
	For(yield func(state State) bool, roots ...State) Trace
}

type BreadthFirstSearch struct {
	system  *TransitionSystem
	forward bool
}

func NewBreadthFirstSearch(system *TransitionSystem, forward bool) BreadthFirstSearch {
	return BreadthFirstSearch{
		system:  system,
		forward: forward,
	}
}

func (search BreadthFirstSearch) successors() func(node State) []State {
	if search.forward {
		return search.system.Outgoing
	}
	return search.system.Ingoing
}

func (search BreadthFirstSearch) For(yield func(state State) bool, roots ...State) Trace {
	var terminal structures.LinkedNode[State]
	states := NewStateSet(roots...)
	algorithms.BFS(
		search.successors(),
		func(state State) bool {
			return states.Contains(state)
		},
		func(node structures.LinkedNode[State]) bool {
			states.Insert(node.Data)
			stop := yield(node.Data)
			if stop {
				terminal = node
			}
			return stop
		},
		roots...,
	)
	return terminal.Array()
}

type DepthFirstSearch struct {
	system  *TransitionSystem
	forward bool
}

func NewDepthFirstSearch(system *TransitionSystem, forward bool) DepthFirstSearch {
	return DepthFirstSearch{
		system:  system,
		forward: forward,
	}
}

func (search DepthFirstSearch) successors() func(node State) []State {
	if search.forward {
		return search.system.Outgoing
	}
	return search.system.Ingoing
}

func (search DepthFirstSearch) For(yield func(state State) bool, roots ...State) Trace {
	var terminal structures.LinkedNode[State]
	states := NewStateSet(roots...)
	algorithms.DFS(
		search.successors(),
		func(state State) bool {
			return states.Contains(state)
		},
		func(node structures.LinkedNode[State]) bool {
			states.Insert(node.Data)
			stop := yield(node.Data)
			if stop {
				terminal = node
			}
			return stop
		},
		roots...,
	)
	return terminal.Array()
}
