package automata

import (
	"github.com/Brandhoej/gobion/pkg/algorithms"
	"github.com/Brandhoej/gobion/pkg/structures"
)

type SearchStrategy interface {
	For(yield func(state State) bool, roots ...State) Trace
}

type BreadthFirstSearch struct {
	system *TransitionSystem
	solver *ConstraintSolver
}

func NewBreadthFirstSearch(system *TransitionSystem, solver *ConstraintSolver) BreadthFirstSearch {
	return BreadthFirstSearch{
		system: system,
		solver: solver,
	}
}

func (search BreadthFirstSearch) For(yield func(state State) bool, roots ...State) Trace {
	var terminal structures.LinkedNode[State]
	states := NewStateSet()
	states.Insert(search.solver, roots...)
	algorithms.BFS(
		search.system.Outgoing,
		func(state State) bool {
			return states.Contains(state, search.solver)
		},
		func(node structures.LinkedNode[State]) bool {
			states.Insert(search.solver, node.Data)
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
	system *TransitionSystem
	solver *ConstraintSolver
}

func NewDepthFirstSearch(system *TransitionSystem, solver *ConstraintSolver) DepthFirstSearch {
	return DepthFirstSearch{
		system: system,
		solver: solver,
	}
}

func (search DepthFirstSearch) For(yield func(state State) bool, roots ...State) Trace {
	var terminal structures.LinkedNode[State]
	states := NewStateSet()
	states.Insert(search.solver, roots...)
	algorithms.DFS(
		search.system.Outgoing,
		func(state State) bool {
			return states.Contains(state, search.solver)
		},
		func(node structures.LinkedNode[State]) bool {
			states.Insert(search.solver, node.Data)
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
