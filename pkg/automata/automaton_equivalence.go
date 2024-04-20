package automata

type AutomatonEquivalence interface {
	Equivalent(lhs, rhs Automaton) bool
}

type StructuralAutomatonEquivalence struct {}

func NewStructuralAutomatonEquivalence() StructuralAutomatonEquivalence {
	return StructuralAutomatonEquivalence{}
}

func (equivalence StructuralAutomatonEquivalence) Equivalent(lhs, rhs Automaton) bool {
	return false
}