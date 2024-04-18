package automata

type Location struct {
	name string
	invariant Invariant
}

func NewLocation(name string, invariant Invariant) Location {
	return Location{
		name: name,
		invariant: invariant,
	}
}

func (location Location) IsEnabled(variables Variables, valuations Valuations) bool {
	return location.invariant.IsSatisfiable(variables, valuations)
}