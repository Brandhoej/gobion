package zones

import (
	"fmt"
	"io"

	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/kelindar/bitmap"
)

// A dense (C x C) matrix representing relations between all clocks "C".
type DBM struct {
	clocks Clock
	data   []Relation
}

// Constructs a DBM with clocks all with the same element filler.
// Often the element filler will be the Zero element.
func NewDBM(clocks Clock, filler Relation) DBM {
	if clocks == 0 {
		panic("DBM require at least one clock that can be the reference clock")
	}

	cardinality := clocks * clocks
	data := make([]Relation, cardinality)
	dbm := DBM{
		clocks: clocks,
		data:   data,
	}

	for row := Reference; row < dbm.clocks; row++ {
		for column := Reference; column < dbm.clocks; column++ {
			if row == column {
				dbm.Constrain(row, column, Zero)
			} else if row == 0 {
				dbm.Constrain(row, column, Zero)
			} else {
				dbm.Constrain(row, column, filler)
			}
		}
	}

	return dbm
}

// Creates a new copy of the DBM.
func (dbm DBM) Copy() DBM {
	data := make([]Relation, dbm.clocks*dbm.clocks)
	copy(data, dbm.data)
	return DBM{
		clocks: dbm.clocks,
		data:   data,
	}
}

// Uses the row-wise indexing and not the layered approach since we have the clock set in the DBM.
// Eg. 3 clocks (including the reference clock) DBM indexing "(row; column)-index":
//
// [(0; 0)-0, (0; 1)-1, (0; 2)-2]
//
// [(1; 0)-3, (1; 1)-4, (1; 2)-5]
//
// [(2; 0)-6, (2; 1)-7, (2; 2)-8]
func (dbm DBM) Index(row, column Clock) uint32 {
	return uint32(row*dbm.clocks + column)
}

// Returns the constraint at the row/column of the DBM.
func (dbm DBM) Constraint(row, column Clock) Relation {
	return dbm.data[dbm.Index(row, column)]
}

// Sets the constraint of the row/column.
func (dbm DBM) Constrain(row, column Clock, relation Relation) {
	dbm.data[dbm.Index(row, column)] = relation
}

// closes the DBm by applying floyds algorithm on it.
func (dbm DBM) Clock() {
	for k := Reference; k < dbm.clocks; k++ {
		for i := Reference; i < dbm.clocks; i++ {
			for j := Reference; j < dbm.clocks; j++ {
				pathIKJ := dbm.Constraint(i, k).Add(
					dbm.Constraint(k, j),
				)
				// i -> k -> j is shorter than i -> j
				if pathIKJ < dbm.Constraint(i, j) {
					dbm.Constrain(i, j, pathIKJ)
				}
			}
		}
	}
}

// Computes the relation lhs has to rhs. If subset is true
// then lhs ⊆ rhs and if superset is true then lhs ⊇ rhs.
// Implied is that if both are true then lhs = rhs.
// The from and to clocks represents the part of the DBMs to compare.
// This allows users to specify a subset of either of the DBMs to compare.
func (lhs DBM) Relation(rhs DBM, from, to Clock) (subset bool, superset bool) {
	subset, superset = true, true
	for row := from; row < to; row++ {
		for column := from; column < to; column++ {
			if !subset && !superset {
				break
			}

			lhsConstraint := lhs.Constraint(row, column)
			rhsConstraint := rhs.Constraint(row, column)

			subset = subset && (lhsConstraint <= rhsConstraint)
			superset = superset && (lhsConstraint >= rhsConstraint)
		}
	}

	return subset, superset
}

// Makes the lhs the intersection of the clocks [from, to] and returns true
// if there actually is a valid intersection. Otherwise, false.
func (lhs DBM) Intersection(rhs DBM, from, to Clock) bool {
	for row := from; row < to; row++ {
		for column := from; column < to; column++ {
			if lhs.Constraint(row, column) > rhs.Constraint(row, column) {
				lhs.Constrain(row, column, rhs.Constraint(row, column))
				if rhs.Constraint(row, column).Negation() >= lhs.Constraint(column, row) {
					return false
				}
			}
		}
	}

	return true
}

// Makes the lhs the convex union of the clocks [from, to] and returns true
// if there actually is a valid intersection. Otherwise, false.
func (lhs DBM) ConvexUnion(rhs DBM, from, to Clock) {
	for row := from; row < to; row++ {
		for column := from; column < to; column++ {
			if lhs.Constraint(row, column) < rhs.Constraint(row, column) {
				lhs.Constrain(row, column, rhs.Constraint(row, column))
			}
		}
	}
}

// Sets the constraint and ensure that the DBM is closed.
// However, it does not check if diagonals are invalid.
func (dbm DBM) ConstrainAndClose(row, column Clock, element Relation) {
	// Only if we are reducing the amount of valuations do we actually constrain it.
	if dbm.Constraint(row, column) > element {
		dbm.Constrain(row, column, element)
		// If the negation is greater than or equal to column -> row.
		// If this is the case then it is implied that there exists a
		// timing constraint such that time can't progress from.
		// The check essentially verifies whether the new constraint
		// introduces a contradiction that implies clocks i and j should
		// both be ahead or behind of each other in terms of time,
		// which is not feasible.
		if element.Negation() >= dbm.Constraint(column, row) {
			dbm.Empty()
		}

		dbm.CloseRowColumn(row, column)
	}
}

// Closes the DBM by recomputing all affected shortest paths.
// This is down by potentially updating the following three paths:
//   - row -> column -> i => row -> i
//   - i -> row -> column => i -> column
//   - i -> column -> j => i -> j
//
// Where "i" and "j" are both clocks (Including the reference clock).
func (dbm DBM) CloseRowColumn(row, column Clock) {
	if dbm.clocks <= 2 {
		return
	}

	pathRC := dbm.Constraint(row, column)

	// Computes the path from "row" to a clock "i" through "column".
	// row -> column -> i => row -> i
	for i := Reference; i < dbm.clocks; i++ {
		pathCI := dbm.Constraint(column, i)
		if pathCI.IsInfinity() {
			continue
		}

		// Check if "row -> c" is shorter through "column".
		pathRC := pathRC.Add(pathCI)
		if dbm.Constraint(row, i) > pathRC {
			dbm.Constrain(row, i, pathRC)
		}
	}

	// Computes the path from a clock "i" to "column" through "row".
	// i -> row -> column => i -> column
	// i -> column -> j => i -> j
	for i := Reference; i < dbm.clocks; i++ {
		// i -> r
		pathIR := dbm.Constraint(i, row)
		if pathIR.IsInfinity() {
			continue
		}

		// i -> row -> column => i -> column
		pathIC := pathIR.Add(pathRC)
		if dbm.Constraint(i, column) <= pathIC {
			continue
		}
		dbm.Constrain(i, column, pathIC)

		for j := Reference; j < dbm.clocks; j++ {
			// column -> j
			pathCJ := dbm.Constraint(column, j)
			if pathCJ.IsInfinity() {
				continue
			}

			// i -> column -> j => i -> j
			pathIJ := pathIC.Add(pathCJ)
			if dbm.Constraint(i, j) > pathIJ {
				dbm.Constrain(i, j, pathIJ)
			}
		}
	}
}

// Applies Floyds algorithm to create a transitive closure.
func (dbm DBM) Close() {
	for k := Reference; k < dbm.clocks; k++ {
		for i := Reference; i < dbm.clocks; i++ {
			if i == k {
				continue
			}

			for j := Reference; j < dbm.clocks; j++ {
				pathIK := dbm.Constraint(i, k)
				if pathIK.IsInfinity() {
					continue
				}

				pathKJ := dbm.Constraint(k, j)
				if pathKJ.IsInfinity() {
					continue
				}

				pathIJ := dbm.Constraint(i, j)
				pathIKJ := pathIK.Add(pathKJ)
				if pathIJ > pathIKJ {
					dbm.Constrain(i, j, pathIKJ)
				}

				if dbm.Constraint(i, j) < Zero {
					dbm.Empty()
				}
			}
		}
	}
}

// Returns the non-redudant edges.
func (dbm DBM) Reduction() graph.Edges[Clock, Constraint] {
	var redundant bitmap.Bitmap = make([]uint64, 0, 1)
	for k := Clock(1); k < dbm.clocks; k++ {
		for i := Clock(1); i < dbm.clocks; i++ {
			if i == k {
				continue
			}

			for j := Clock(1); j < dbm.clocks; j++ {
				pathIK := dbm.Constraint(i, k)
				if pathIK.IsInfinity() {
					continue
				}

				pathKJ := dbm.Constraint(k, j)
				if pathKJ.IsInfinity() {
					continue
				}

				pathIJ := dbm.Constraint(i, j)
				pathIKJ := pathIK.Add(pathKJ)
				if pathIJ <= pathIKJ {
					redundant.Set(uint32(dbm.Index(i, j)))
				}
			}
		}
	}

	edges := make(graph.EdgeSlice[Clock, Constraint], 0)
	for row := Clock(1); row < dbm.clocks; row++ {
		for column := Clock(1); column < dbm.clocks; column++ {
			if !redundant.Contains(dbm.Index(row, column)) {
				constraint := NewConstraint(row, column, dbm.Constraint(row, column))
				edges = append(edges, constraint)
			}
		}
	}

	return &edges
}

// Applies Floyds algorithm to see if the DBM is closed.
func (dbm DBM) IsClosed() bool {
	for i := Reference; i < dbm.clocks; i++ {
		for row := Reference; row < dbm.clocks; row++ {
			for column := Reference; column < dbm.clocks; column++ {
				pathRowI := dbm.Constraint(row, column)
				if pathRowI.IsInfinity() {
					continue
				}

				pathIColumn := dbm.Constraint(i, column)
				if pathIColumn.IsInfinity() {
					continue
				}

				pathRowColumn := dbm.Constraint(row, column)
				pathRowIColumn := pathRowI.Add(pathIColumn)
				if pathRowColumn > pathRowIColumn {
					return false
				}
			}
		}
	}

	return true
}

// checks if the constraints of both DBMs are the same in the range of clocks in "from" and "to".
func (lhs DBM) Equals(rhs DBM, from, to Clock) bool {
	for row := from; row < to; row++ {
		for column := from; column < to; column++ {
			if lhs.Constraint(row, column) != rhs.Constraint(row, column) {
				return false
			}
		}
	}
	return true
}

// Sets the reference clock constraint to negative.
func (dbm DBM) Empty() {
	dbm.SetDiagonal(Reference, NewRelation(-1, 0))
}

// Checks if the DBM is empty.
func (dbm DBM) IsConsistent() bool {
	if dbm.Diagonal(Reference) < Zero {
		return false
	}

	for clock := Reference; clock < dbm.clocks; clock++ {
		if dbm.Upper(clock) < dbm.Lower(clock) ||
			dbm.Diagonal(clock) < Zero {
			dbm.Empty()
			return false
		}
	}

	return true
}

// Returns the upper bound relation that is "clock - 0 ~ n"
// Where "0" is the reference clock making the potential relations:
// "clock < n" and "clock ≤ n" thereby the upper bound of the clock.
func (dbm DBM) Upper(clock Clock) Relation {
	return dbm.Constraint(clock, Reference)
}

// Sets the upper bound on the clock which is the relation "clock - 0 ~ n".
// Where "0" is the reference clock making the potential relations:
// "clock < n" and "clock ≤ n" thereby the upper bound of the clock.
func (dbm DBM) SetUpper(clock Clock, element Relation) {
	// clock - 0 ~ n, where ~ ∈ {<, ≤}.
	// In other words "clock ~ n"
	// Where "n" is the upper bound on "clock" given "~".
	dbm.Constrain(clock, Reference, element)
}

// Returns the lower bound relation that is "0 - clock ~ n"
// Where "0" is the reference clock making the potential relations:
// "-clock < n" and "-clock ≤ n" thereby the upper bound of the clock.
func (dbm DBM) Lower(clock Clock) Relation {
	return dbm.Constraint(Reference, clock)
}

// Sets the lower bound on the clock which is the relation "0 - clock ~ n"
// Where "0" is the reference clock making the potential relations:
// "-clock < n" and "-clock ≤ n" thereby the upper bound of the clock.
func (dbm DBM) SetLower(clock Clock, element Relation) {
	// 0 - clock ~ n, where ~ ∈ {<, ≤}.
	// In other words "-clock ~ n"
	// Where "n" is the lower bound on "clock" given "~".
	dbm.Constrain(Reference, clock, element)
}

// Returns the diagonal relation of the clock.
func (dbm DBM) Diagonal(clock Clock) Relation {
	return dbm.Constraint(clock, clock)
}

// Sets the diagonal of the clock. For it to be consistent it must be Zero.
func (dbm DBM) SetDiagonal(clock Clock, relation Relation) {
	dbm.Constrain(clock, clock, relation)
}

// The up operation computes the strongest postcondition of a zone with respect to delay.
// Afterwards the DBM contains the clock assignments that can be reached from by delay.
// up(D) = {u + d | u ∈ D, d ∈ ℝ+}.
// This operation preserves the canonical form thereby applying it on a canonical DBM
// will result in a new canonical DBM.
func (dbm DBM) Up() {
	// remove the upper bounds on all individual clocks.
	// This is done by setting all elements At(i, 0) to ∞.
	for i := Clock(1); i < dbm.clocks; i++ {
		// We can delay infinitely much relative to the reference clock.
		dbm.SetUpper(i, Infinity)
	}
}

// In contrast to Up, Down computes the weakest precondition of the DBM withrespect to delay.
// down(D) = {u | u + d ∈ D, d ∈ ℝ+} such that the set of clock assignments that can reach D
// by some delay d. Algorithmically, it is computed by setting the lower bound on all individual
// clocks to (0, ≤).
func (dbm DBM) Down() {
	// We are using i based indexing and therefor being in the same i
	// will be faster than switching each iteration.
	for i := Clock(1); i < dbm.clocks; i++ {
		// Only if the lower bound is not already lowered do we lower it.
		// For DBMs where diagonals are valid i.e., they are non-negative
		// then we can assume that only when it is greater than (0, ≤)
		// should it be lowered even further.
		if dbm.Lower(i) != Zero {
			dbm.SetLower(i, Zero)

			for j := Clock(1); j < dbm.clocks; j++ {
				// We check if the constraint between the row and column clocks
				// are stricter than the the lower bound of the column. This ensures that
				// the lower bound on the column is lowest possible considering all
				// difference constraints. This is provided that there even is a finite
				// lower bound on the column. If the lower bound were infinity, it would
				// imply that there is no lower bound for the constraint between clock.
				if dbm.Constraint(j, i) < dbm.Lower(i) {
					// The difference bound between the row and coulmn was stricter than
					// that of the reference clock.
					dbm.SetLower(i, dbm.Constraint(j, i))
				}
			}
		}
	}
}

// Returns true if the relation over the two clocks lies within the DBM.
func (dbm DBM) Satisfies(row, column Clock, relation Relation) bool {
	return dbm.Constraint(row, column) >= relation && dbm.Constraint(column, row) <= relation.Negation()
}

// Returns true if all clocks' upper bound is infinity.
func (dbm DBM) CanDelayIndefinitely() bool {
	for clock := Clock(1); clock < dbm.clocks; clock++ {
		if !dbm.Upper(clock).IsInfinity() {
			return false
		}
	}
	return true
}

// Removes all constraints on a given clock, i.e., the clock may take any positive value.
// This is expressed as {u[x=d] | u ∈ D, d ∈ ℝ+}.
func (dbm DBM) Free(clocks ...Clock) {
	for _, clock := range clocks {
		for dimension := Reference; dimension < dbm.clocks; dimension++ {
			if dimension != clock {
				dbm.Constrain(clock, dimension, Infinity)
				dbm.Constrain(dimension, clock, dbm.Upper(dimension))
			}
		}
	}
}

// Sets the clock to be assigned to its limit. This is expressed as {u[x=m] | u ∈ D}.
func (dbm DBM) Reset(clock Clock, limit int) {
	positive := NewRelation(limit, LessThanEqual)
	negative := NewRelation(-limit, LessThanEqual)
	for dimension := Reference; dimension < dbm.clocks; dimension++ {
		dbm.Constrain(clock, dimension, positive.Add(dbm.Lower(dimension)))
		dbm.Constrain(dimension, clock, dbm.Upper(dimension).Add(negative))
	}
}

// Sets the lhs to be equal to the rhs. This is expressed as {u[x=u(y)] | u ∈ D, x ∈ D}
func (dbm DBM) Assign(lhs, rhs Clock) {
	for dimension := Reference; dimension < dbm.clocks; dimension++ {
		if dimension != lhs {
			dbm.Constrain(lhs, dimension, dbm.Constraint(rhs, dimension))
			dbm.Constrain(dimension, lhs, dbm.Constraint(dimension, rhs))
		}
	}

	dbm.Constrain(lhs, rhs, Zero)
	dbm.Constrain(rhs, lhs, Zero)
}

// Compound addition assignment of the clock "clock := clock + limit".
func (dbm DBM) Shift(clock Clock, limit int) {
	pos := NewRelation(limit, LessThanEqual)
	neg := NewRelation(-limit, LessThanEqual)
	for i := Reference; i < dbm.clocks; i++ {
		// For all, except the loop, we extend the bounds depending on whether the
		// constraint describes if the clock is in front of or behind "i".
		if i != clock {
			// "i" is in front of "clock" and therefore we add "limit".
			from := dbm.Constraint(clock, i)
			dbm.Constrain(clock, i, from.Add(pos))
			// "clock" is in front of "i" and therefore we subtract "limit".
			to := dbm.Constraint(i, clock)
			dbm.Constrain(i, clock, to.Add(neg))
		}
	}

	// It might be that the lower bound is inconsistent and should be clamped.
	if dbm.Lower(clock) > Zero {
		dbm.SetLower(clock, Zero)
	}
	// It might be that the upper bound is inconsistent and should be clamped.
	if dbm.Upper(clock) < Zero {
		dbm.SetUpper(clock, Zero)
	}
}

// Removes all upper bounds higher than the maximal constants and lower
// all lower bounds higher than the maximal constants down to the maximal constants.
//
// If checking safety properties then this can only be used if there are NO difference constraints.
func (dbm DBM) Norm(maximums ...int) {
	for i := Clock(1); i < dbm.clocks; i++ {
		pos := NewRelation(maximums[i-1], LessThanEqual)
		neg := NewRelation(maximums[i-1], LessThan)

		if dbm.Upper(i) > pos {
			dbm.SetUpper(i, Infinity)
		}

		if dbm.Lower(i) > neg {
			dbm.SetLower(i, neg)
		}

		for j := Clock(1); j < dbm.clocks; j++ {
			if i == j {
				continue
			}

			constraint := dbm.Constraint(i, j)
			if constraint.IsInfinity() {
				continue
			}

			if constraint > pos {
				dbm.Constrain(i, j, Infinity)
			} else if constraint < neg {
				dbm.Constrain(i, j, neg)
			}
		}
	}

	dbm.Close()
}

// Writes the array where each element is the tuple or infinity defining the
// relation between the row and column clocks.
func (dbm DBM) WriteMatrix(writer io.Writer) {
	for row := Reference; row < dbm.clocks; row++ {
		io.WriteString(writer, "[")
		for column := Reference; column < dbm.clocks; column++ {
			if column > Reference {
				io.WriteString(writer, ", ")
			}
			io.WriteString(writer, dbm.Constraint(row, column).String())
		}
		io.WriteString(writer, "]\n")
	}
}

// Writes to the writer a all conjunctions representing the DBM.
// The labels are used for each of the clocks, the labels does not include the reference.
// This can be used in GeoGebra Classic to draw the DBM.
func (dbm DBM) WriteConjunctions(writer io.Writer, labels ...string) {
	hasContraints := false
	for row := Reference; row < dbm.clocks; row++ {
		for column := Reference; column < dbm.clocks; column++ {
			if row == Reference && column == Reference {
				continue
			}

			element := dbm.Constraint(row, column)
			if !element.IsInfinity() {
				if hasContraints {
					io.WriteString(writer, " ∧ ")
				}

				limit := element.Limit()
				if row == 0 {
					// Lower bound.
					rhs := labels[column-1]
					io.WriteString(writer, fmt.Sprintf("-%s %s %v", rhs, element.Strictness().String(), limit))
					hasContraints = true
				} else if column == 0 {
					// Upper bound.
					lhs := labels[row-1]
					io.WriteString(writer, fmt.Sprintf("%s %s %v", lhs, element.Strictness().String(), limit))
					hasContraints = true
				} else {
					// Diference relation.
					lhs, rhs := labels[row-1], labels[column-1]
					io.WriteString(writer, fmt.Sprintf("%s - %s %s %v", lhs, rhs, element.Strictness().String(), limit))
					hasContraints = true
				}

			}
		}
	}
}

// Writes the complete graph interpretation of a DBM to the writer. The labels
// are used for each of the clocks, the labels does include the reference.
// This can be used with Graphviz.
func (dbm DBM) WriteGraphvizDigraph(writer io.Writer, distances bool, labels ...string) {
	io.WriteString(writer, "digraph {\n")

	for i := Reference; i < dbm.clocks; i++ {
		io.WriteString(
			writer,
			fmt.Sprintf("%v [label=\"%s\" shape=\"circle\"]\n", i, labels[i]),
		)
	}

	for row := Reference; row < dbm.clocks; row++ {
		for column := Reference; column < dbm.clocks; column++ {
			constraint := dbm.Constraint(row, column)
			label := constraint.String()
			if distances {
				label = fmt.Sprintf("%v", constraint.Limit())
			}
			io.WriteString(
				writer,
				fmt.Sprintf("%v -> %v [label=\"%s\"]\n", row, column, label),
			)
		}
	}

	io.WriteString(writer, "}")
}
