package golang

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
	   	x_0_0  := 0
	   	x_0_1 := x_0 + 1
	   	x_0_2 := x_1 + 2
	*/
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	scope := NewGoGlobalScope(0, context.NewSolver())
	sort := context.IntegerSort()
	identifier := "x"

	// Act
	x0 := scope.Declare(identifier, context.NewInt(0, sort))            // x := 0
	x0Val, _ := scope.Valuation(identifier)                             // 		=>  x_0  := 0
	x1 := scope.Assign(identifier, z3.Add(x0, context.NewInt(1, sort))) // x += 1
	x1Val, _ := scope.Valuation(identifier)                             // 		=>  x_1 := x_0 + 1
	x2 := scope.Assign(identifier, z3.Add(x1, context.NewInt(2, sort))) // x += 2
	x2Val, _ := scope.Valuation(identifier)                             //		=>  x_2 := x_1 + 2

	// Assert
	assert.Equal(t, "x_0_0", x0.String())
	assert.Equal(t, "0", x0Val.String())
	assert.Equal(t, "x_0_1", x1.String())
	assert.Equal(t, "(+ x_0_0 1)", x1Val.String())
	assert.Equal(t, "x_0_2", x2.String())
	assert.Equal(t, "(+ x_0_1 2)", x2Val.String())
}

func TestSingleStaticAssignmentMultipleVariable(t *testing.T) {
	/*
	   	x := 0
	   	x += 1
	   	y := 1
	   	x += 2
	   	y += x + y

	   SSA:
	   	x_0_0 := 0
	   	x_0_1 := x_0 + 1
	   	y_0_0 := 1
	   	x_0_2 := x_1 + 2
	   	y_0_1 := x_2 + y
	*/
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	scope := NewGoGlobalScope(0, context.NewSolver())
	sort := context.IntegerSort()
	xIdent, yIdent := "x", "y"

	// Act
	x0 := scope.Declare(xIdent, context.NewInt(0, sort))            // x := 0
	x0Val, _ := scope.Valuation(xIdent)                             // 		=>  x_0 := 0
	x1 := scope.Assign(xIdent, z3.Add(x0, context.NewInt(1, sort))) // x += 1
	x1Val, _ := scope.Valuation(xIdent)                             // 		=>  x_1 := x_0 + 1
	y := scope.Declare(yIdent, context.NewInt(1, sort))             // y := 1
	yVal, _ := scope.Valuation(yIdent)                              // 		=>  y_0 := 1
	x2 := scope.Assign(xIdent, z3.Add(x1, context.NewInt(2, sort))) // x += 2
	x2Val, _ := scope.Valuation(xIdent)                             //		=>  x_2 := x_1 + 2
	y1 := scope.Assign(yIdent, z3.Add(x2, y))                       // y += x + y
	y1Val, _ := scope.Valuation(yIdent)                             // 		=>  y_1 := x_2 + y

	// Assert
	assert.Equal(t, "x_0_0", x0.String())
	assert.Equal(t, "0", x0Val.String())
	assert.Equal(t, "x_0_1", x1.String())
	assert.Equal(t, "(+ x_0_0 1)", x1Val.String())
	assert.Equal(t, "y_0_0", y.String())
	assert.Equal(t, "1", yVal.String())
	assert.Equal(t, "x_0_2", x2.String())
	assert.Equal(t, "(+ x_0_1 2)", x2Val.String())
	assert.Equal(t, "y_0_1", y1.String())
	assert.Equal(t, "(+ x_0_2 y_0_0)", y1Val.String())
}

func TestVariableShadowing(t *testing.T) {
	/*
	   	x := 0
	   	x = 1
	   	{
	   		x := 0
	   		x = 2
	   	}

	   SSA:
	   	x_0_0 := 0
	   	x_0_1 := 1
	   	x_1_0 := 0
	   	x_1_1 := 2
	*/
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	outer := NewGoGlobalScope(0, context.NewSolver())
	inner := outer.Child(outer.id+1)
	sort := context.IntegerSort()
	xIdent := "x"

	// Act
	x0 := outer.Declare(xIdent, context.NewInt(0, sort))
	x0Val, _ := outer.Valuation(xIdent)
	x1 := outer.Assign(xIdent, context.NewInt(1, sort))
	x1Val, _ := outer.Valuation(xIdent)
	x10 := inner.Declare(xIdent, context.NewInt(0, sort))
	x10Val, _ := inner.Valuation(xIdent)
	x11 := inner.Assign(xIdent, context.NewInt(2, sort))
	x11Val, _ := inner.Valuation(xIdent)

	// Assert
	assert.Equal(t, "x_0_0", x0.String())
	assert.Equal(t, "0", x0Val.String())
	assert.Equal(t, "x_0_1", x1.String())
	assert.Equal(t, "1", x1Val.String())

	assert.Equal(t, "x_1_0", x10.String())
	assert.Equal(t, "0", x10Val.String())
	assert.Equal(t, "x_1_1", x11.String())
	assert.Equal(t, "2", x11Val.String())
}

func TestIsLocal(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	outer := NewGoGlobalScope(0, context.NewSolver())
	inner := outer.Child(outer.id+1)
	sort := context.IntegerSort()
	xIdent, yIdent, zIdent := "x", "y", "z"

	// Act
	outer.Define(xIdent, sort)
	outer.Define(zIdent, sort)
	inner.Define(yIdent, sort)
	inner.Define(zIdent, sort)

	// Assert
	assert.True(t, outer.IsAccessible(xIdent), "x should be accessible in the outer scope")
	assert.True(t, inner.IsAccessible(xIdent), "x should be accessible in the inner scope")
	assert.True(t, outer.IsLocal(xIdent), "x should be local in the outer scope")
	assert.False(t, inner.IsLocal(xIdent), "x should not be local in the inner scope")

	assert.True(t, inner.IsLocal(yIdent), "y should be local in the inner scope")
	assert.False(t, outer.IsLocal(yIdent), "y should not be local in the outer scope")
	assert.False(t, outer.IsAccessible(yIdent), "y should not be accessible in the outer scope")
	assert.True(t, inner.IsAccessible(yIdent), "y should be accessible in the inner scope")

	assert.True(t, outer.IsLocal(zIdent), "z should be local in the outer scope")
	assert.True(t, outer.IsAccessible(zIdent), "z should be accessible in the outer scope")
	assert.True(t, inner.IsLocal(zIdent), "z in inner scope should shadow z in outer scope")
	assert.True(t, inner.IsAccessible(zIdent), "z should be accessible in the inner scope")
}

func TestAssign(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	solver := context.NewSolver()
	outer := NewGoGlobalScope(0, context.NewSolver())
	inner := outer.Child(outer.id+1)
	sort := context.IntegerSort()
	one := context.NewInt(1, sort)
	xIdent := "x"

	// Act
	outer.Declare(xIdent, one)
	outerValuation, outerExists := outer.Valuation(xIdent)
	innerValuation, innerExists := inner.Valuation(xIdent)

	// Assert
	assert.True(t, outerExists)
	assert.True(t, innerExists)
	if model := solver.Prove(z3.Eq(outerValuation, one)); model != nil {
		t.Error(z3.Eq(outerValuation, one).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(innerValuation, one)); model != nil {
		t.Error(z3.Eq(innerValuation, one).String(), "has solution", model.String())
	}
}

func TestVariable(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	outer := NewGoGlobalScope(0, context.NewSolver())
	inner := outer.Child(outer.id+1)
	sort := context.IntegerSort()
	xIdent, yIdent, zIdent := "x", "y", "z"

	// Act
	x := outer.Define(xIdent, sort)
	zo := outer.Define(zIdent, sort)
	y := inner.Define(yIdent, sort)
	zi := inner.Define(zIdent, sort)

	xVar, xVarExists := outer.Variable(xIdent)
	zoVar, zoVarExists := outer.Variable(zIdent)
	yVar, yVarExists := inner.Variable(yIdent)
	ziVar, ziVarExists := inner.Variable(zIdent)

	// Assert
	assert.True(t, xVarExists, "x variable should exist in the outer scope")
	assert.True(t, zoVarExists, "z variable should exist in the outer scope")
	assert.True(t, yVarExists, "y variable should exist in the inner scope")
	assert.True(t, ziVarExists, "z variable in the inner scope should shadow z in the outer scope")
	assert.Equal(t, x.String(), xVar.String(), "x declaration and variable should be the same")
	assert.Equal(t, zo.String(), zoVar.String(), "z declaration and variable should be the same")
	assert.Equal(t, y.String(), yVar.String(), "y declaration and variable should be the same")
	assert.Equal(t, zi.String(), ziVar.String(), "z declaration and variable should be the same")
}

func TestValuation(t *testing.T) {
	// Arrange
	config := z3.NewConfig()
	context := z3.NewContext(config)
	solver := context.NewSolver()
	outer := NewGoGlobalScope(0, context.NewSolver())
	inner := outer.Child(outer.id+1)
	sort := context.IntegerSort()
	one := context.NewInt(1, sort)
	two := context.NewInt(2, sort)
	three := context.NewInt(3, sort)
	xIdent := "x"

	// Act
	x1 := outer.Declare(xIdent, one)
	outer1, outer1Exists := outer.Valuation(xIdent)
	inner1, inner1Exists := inner.Valuation(xIdent)
	x2 := inner.Assign(xIdent, two)
	outer2, outer2Exists := outer.Valuation(xIdent)
	inner2, inner2Exists := inner.Valuation(xIdent)
	x3 := inner.Declare(xIdent, three)
	outer3, outer3Exists := outer.Valuation(xIdent)
	inner3, inner3Exists := inner.Valuation(xIdent)

	// Assert
	assert.True(t, outer1Exists)
	assert.True(t, inner1Exists)
	assert.True(t, outer2Exists)
	assert.True(t, inner2Exists)
	assert.True(t, outer3Exists)
	assert.True(t, inner3Exists)
	assert.NotNil(t, x1, "declaration of x in the outer scope should return the variable x")
	assert.NotNil(t, x2, "assignemt of x in the inner scope should return the variable x from the outer scope")
	assert.NotNil(t, x3, "declaration of x in the inner scope should shadow x from the outer scope")
	assert.Equal(t, "x_0_0", x1.String())
	assert.Equal(t, "x_0_1", x2.String())
	assert.Equal(t, "x_1_0", x3.String())
	if model := solver.Prove(z3.Eq(outer1, one)); model != nil {
		t.Error(z3.Eq(outer1, one).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(inner1, one)); model != nil {
		t.Error(z3.Eq(inner1, one).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(outer2, two)); model != nil {
		t.Error(z3.Eq(outer2, two).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(inner2, two)); model != nil {
		t.Error(z3.Eq(inner2, two).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(outer3, two)); model != nil {
		t.Error(z3.Eq(outer3, two).String(), "has solution", model.String())
	}
	if model := solver.Prove(z3.Eq(inner3, three)); model != nil {
		t.Error(z3.Eq(inner3, two).String(), "has solution", model.String())
	}
}
