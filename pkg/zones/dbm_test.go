package zones

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fig11() DBM {
	dbm := NewDBM(Clock(1+2), Infinity)
	dbm.SetLower(Clock(1), NewRelation(-1, Strict))
	dbm.SetUpper(Clock(1), NewRelation(3, Strict))
	dbm.SetLower(Clock(2), NewRelation(-2, Strict))
	dbm.SetUpper(Clock(2), NewRelation(3, Strict))
	return dbm
}

func Test_Fig11Copy(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	copy := dbm.Copy()

	// Assert
	assert.True(t, dbm.Equals(copy, Reference, Clock(3)))
}

func Test_Fig11Up(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Up()

	// Assert
	assert.Equal(t, dbm.Upper(Reference), Zero)
	assert.Equal(t, dbm.Upper(Clock(1)), Infinity)
	assert.Equal(t, dbm.Upper(Clock(2)), Infinity)
}

func Test_Fig11Down(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Down()

	// Assert
	assert.Equal(t, dbm.Lower(Reference), Zero)
	assert.Equal(t, dbm.Lower(Clock(1)), Zero)
	assert.Equal(t, dbm.Lower(Clock(2)), Zero)
}

func Test_Fig11Free(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Free(Clock(2))

	// Assert
	assert.Equal(t, dbm.Lower(Reference), Zero)
	assert.Equal(t, dbm.Lower(Clock(2)), Zero)
	assert.Equal(t, dbm.Upper(Clock(2)), Infinity)
}

func Test_Fig11Reset(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Reset(Clock(1), 2)

	// Assert
	assert.Equal(t, dbm.Upper(Clock(1)), NewRelation(2, Weak))
	assert.Equal(t, dbm.Lower(Clock(1)), NewRelation(-2, Weak))
}

func Test_Fig11Assign(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Assign(Clock(1), Clock(2))

	// Assert
	assert.Equal(t, dbm.Lower(Clock(1)), NewRelation(-2, Strict))
	assert.Equal(t, dbm.Constraint(Clock(1), Clock(2)), Zero)
	assert.Equal(t, dbm.Constraint(Clock(2), Clock(1)), Zero)
}

func Test_Fig11Constrain(t *testing.T) {
	// Arrange
	dbm := fig11()
	constraint := NewConstraint(
		Clock(1), Reference, NewRelation(2, Weak),
	)

	// Act
	dbm.ConstrainAndClose(constraint.i, constraint.j, constraint.relation)

	// Assert
	assert.Equal(t, dbm.Constraint(constraint.i, constraint.j), constraint.relation)
}

func Test_Fig11NormK(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Norm(2, 1)

	// Assert
	var buffer bytes.Buffer
	dbm.WriteConjunctions(&buffer, "x", "y")
	fmt.Println(buffer.String())
	t.Fail()
}

func Test_Fig11Shift(t *testing.T) {
	// Arrange
	dbm := fig11()

	// Act
	dbm.Shift(Clock(2), 1)

	// Assert
	var buffer bytes.Buffer
	dbm.WriteConjunctions(&buffer, "x", "y")
	fmt.Println(buffer.String())
	t.Fail()
}

func Test_foo(t *testing.T) {
	// Arrange
	dbm := NewDBM(Clock(1+2), Infinity)
	dbm.SetLower(Clock(1), NewRelation(-3, Strict))
	dbm.SetUpper(Clock(1), NewRelation(5, Strict))
	dbm.SetLower(Clock(2), NewRelation(-1, Strict))
	dbm.SetUpper(Clock(2), NewRelation(2, Strict))

	// Act
	var buffer bytes.Buffer
	dbm.WriteConjunctions(&buffer, "x", "y")

	// Assert
	fmt.Println(buffer.String())
	t.Fail()
}

func Test_bar(t *testing.T) {
	// Arrange
	dbm := NewDBM(Clock(1+2), Zero)

	// Act
	var buffer bytes.Buffer
	dbm.WriteConjunctions(&buffer, "x", "y")

	// Assert
	fmt.Println(buffer.String())
	t.Fail()
}
