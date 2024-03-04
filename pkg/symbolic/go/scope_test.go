package symbolic

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/stretchr/testify/assert"
)

func TestSingleStaticAssignmentSingleVariable(t *testing.T) {
/*
	x := 0
	x += 1
	x += 2

SSA:
	x_0  := 0
	x_1 := x_0 + 1
	x_2 := x_1 + 2
*/
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	scope := NewGoGlobalScope()
	sort := context.IntegerSort()
	identifier := "x"

	// Act
	x := scope.Declare(identifier, context.NewInt(0, sort))             // x := 0
	xVal, _ := scope.Valuation(identifier)                              // 		=>  x_0  := 0
	x1 := scope.Assign(identifier, z3.Add(x, context.NewInt(1, sort)))  // x += 1
	x1Val, _ := scope.Valuation(identifier)                             // 		=>  x_1 := x_0 + 1
	x2 := scope.Assign(identifier, z3.Add(x1, context.NewInt(2, sort))) // x += 2
	x2Val, _ := scope.Valuation(identifier)                             //		=>  x_2 := x_1 + 2

	// Assert
	assert.Equal(t, "x_0", x.String())
	assert.Equal(t, "0", xVal.String())
	assert.Equal(t, "x_1", x1.String())
	assert.Equal(t, "(+ x_0 1)", x1Val.String())
	assert.Equal(t, "x_2", x2.String())
	assert.Equal(t, "(+ x_1 2)", x2Val.String())
}

func TestSingleStaticAssignmentMultipleVariable(t *testing.T) {
/*
	x := 0
	x += 1
	y := 1
	x += 2
	y += x + y

SSA:
	x_0 := 0
	x_1 := x_0 + 1
	y_0 := 1
	x_2 := x_1 + 2
	y_1 := x_2 + y
*/
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	scope := NewGoGlobalScope()
	sort := context.IntegerSort()
	xIdent, yIdent := "x", "y"

	// Act
	x := scope.Declare(xIdent, context.NewInt(0, sort))             // x := 0
	xVal, _ := scope.Valuation(xIdent)                              // 		=>  x_0 := 0
	x1 := scope.Assign(xIdent, z3.Add(x, context.NewInt(1, sort)))  // x += 1
	x1Val, _ := scope.Valuation(xIdent)                             // 		=>  x_1 := x_0 + 1
	y := scope.Declare(yIdent, context.NewInt(1, sort))             // y := 1
	yVal, _ := scope.Valuation(yIdent)                              // 		=>  y_0 := 1
	x2 := scope.Assign(xIdent, z3.Add(x1, context.NewInt(2, sort))) // x += 2
	x2Val, _ := scope.Valuation(xIdent)                             //		=>  x_2 := x_1 + 2
	y1 := scope.Assign(yIdent, z3.Add(x2, y))                       // y += x + y
	y1Val, _ := scope.Valuation(yIdent)                             // 		=>  y_1 := x_2 + y

	// Assert
	assert.Equal(t, "x_0", x.String())
	assert.Equal(t, "0", xVal.String())
	assert.Equal(t, "x_1", x1.String())
	assert.Equal(t, "(+ x_0 1)", x1Val.String())
	assert.Equal(t, "y_0", y.String())
	assert.Equal(t, "1", yVal.String())
	assert.Equal(t, "x_2", x2.String())
	assert.Equal(t, "(+ x_1 2)", x2Val.String())
	assert.Equal(t, "y_1", y1.String())
	assert.Equal(t, "(+ x_2 y_0)", y1Val.String())
}

func TestVariableShadowing(t *testing.T) {
/*
	x := 0
	x += 1
	{
		x := 1
		x = 2
	}

SSA:
	x  := 0
	x1 := x + 1
	x10 :=
*/
	// Arrange
	// Act
	// Assert
}