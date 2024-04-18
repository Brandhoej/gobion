package z3

import "testing"

func Test_Substitute(t *testing.T) {
	// Arrange
	context := NewContext(NewConfig())
	one := context.NewInt(1, context.IntegerSort())
	two := context.NewInt(2, context.IntegerSort())
	three := context.NewInt(3, context.IntegerSort())
	x := context.NewConstant(WithInt(0), context.IntegerSort())
	y := context.NewConstant(WithInt(1), context.IntegerSort())
	formula := Multiply(x, y, three)

	// Act
	sub := formula.Substitute(
		[]*AST{x, y}, []*AST{two, one},
	)
	t.Log(formula.String())
	t.Log(sub.String())
	t.FailNow()

	// Assert
}
