package zones

import (
	"fmt"
	"math"
)

const (
	// These describe relations as uint but are not the Relation type.
	// The reason for this is that we can omit type conversions from Relation to uint then.
	Strict = 0
	Weak   = 1

	// Represents a DBM element infinity (unbounded) relation between two clocks.
	Infinity = Relation(math.MinInt)
	// Represents a DBM element Zero (0, ≤) relation between two clocks.
	Zero = Relation(Weak)

	// The clock (index) of the reference clock (0 clock) within the DBM.
	Reference = Clock(0)
)

// An index of the clock just like a symbol is for an identifier.
type Clock uint

// A 1 bit representation of either a strict or weak relation (< or ≤).
// A uint is used such that type conversion is omitted.
type Strictness uint

// Returns a string for the relation which is either "<" or "≤".
func (strictness Strictness) String() string {
	if strictness == Strict {
		return "<"
	}
	if strictness == Weak {
		return "≤"
	}
	panic("Unknown relation")
}

// A machine length element optimized for caching which represents a strict
// or weak relation between two clocks. This encoding uses the least significant
// bit to represent the relation and the other bits as the limit.
// INT: [limit] [1 bit relation]
type Relation int

// Constructs an element with a limit and relation.
// If limit is (math.MaxInt << 1) and the relation is < then it is an unbounded infinity relation.
func NewRelation(limit int, strictness Strictness) Relation {
	// The least significant bit is the relation and to the right is the limit.
	// This encoding maintains that strict bounds are smaller than non-strict
	// This allows user to write "e1 < e2" and so on.
	return Relation((limit << 1) | int(strictness))
}

// Constructs an element that represents an unbounded relation between two clocks.
func NewInfinity() Relation {
	return Infinity
}

// Construct an element that represents a zero relation between two clocks.
func NewZero() Relation {
	return Zero
}

// Returns true if the limit is less than infinity - i.e., it is unbounded.
// In other words, the complement is 0.
func (element Relation) IsInfinity() bool {
	return element == Infinity
}

// Returns true if the limit is less than or equal to 0.
// In other words the uint value of the element is 1.
func (element Relation) IsZero() bool {
	return element == Zero
}

// Returns the negation of the element. That is, the relation goes
// from strict to weak and vice versa. Whilst, the limit is negated
// as a regular integer value.
func (element Relation) Negation() Relation {
	if element.IsInfinity() {
		panic("Infinity cannot be negated")
	}

	return 1 - element
}

// Addition of two constraints represented by tuples. The sum constraint
// is compromised such that it satisfies both original constraints (lhs/rhs).
// In other words, it does not allow any behavior that violates either of the
// original constraints. to ensure that the sum captures the intersection of the
// original constraints accurately, we choose the tightest or most restrictive
// relation that satisfies both original constraints. This ensures that the
// resulting constraint is as tight as possible while still being consistent
// with the original constraints. Thereby, if one relation is ≤ (i.e., 1)
// then it is keps over < (i.e., 0). Addition does not handle overflows, and
// therefore yeilds undefined behaviour.
func (lhs Relation) Add(rhs Relation) Relation {
	if lhs.IsInfinity() || rhs.IsInfinity() {
		return NewInfinity()
	}

	// First adding the lhs and rhs increases the limit.
	// Then we ensure the tightest constraint that satisfies both constraints is kept.
	return (lhs + rhs) - ((lhs & Weak) | (rhs & Weak))
}

// Returns the integer bounds, that is "n" of "i-j ~ n".
func (element Relation) Limit() int {
	return int(element) >> 1
}

// Returns the relation, that is "~" of "i-j ~ n".
func (element Relation) Strictness() Strictness {
	return Strictness(element & Weak)
}

// Returns a pretty printed string of the element as a tuple if not infinity.
func (element Relation) String() string {
	if element.IsInfinity() {
		return "∞"
	}
	return fmt.Sprintf("(%v, %s)", element.Limit(), element.Strictness())
}

// Represents a constraint between two clocks i and j: "i-j ~ n".
type Constraint struct {
	i, j     Clock
	relation Relation
}

func (constraint Constraint) Source() Clock {
	return constraint.i
}

func (constraint Constraint) Destination() Clock {
	return constraint.j
}

// Constructs a new constraint with a relation between the two clocks i and j.
func NewConstraint(i, j Clock, relation Relation) Constraint {
	return Constraint{
		i:        i,
		j:        j,
		relation: relation,
	}
}
