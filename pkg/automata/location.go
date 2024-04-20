package automata

type Location struct {
	name      string
	invariant Invariant
}

func NewLocation(name string, invariant Invariant) Location {
	return Location{
		name:      name,
		invariant: invariant,
	}
}

func (location Location) IsEnabled(solver *ConstraintSolver) bool {
	return location.invariant.IsSatisfiable(solver)
}
