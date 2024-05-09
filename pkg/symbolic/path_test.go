package symbolic

import (
	"testing"

	"github.com/Brandhoej/gobion/internal/z3"
	"github.com/stretchr/testify/assert"
)

func TestSequentialPath(t *testing.T) {
	// first ->  second -> third

	// Arrange
	context := z3.NewContext(z3.NewConfig())

	// Act
	first := NewGlobalPath(
		NewLexicalScope(context), context.NewTrue(),
	)
	seocnd := first.Fork(context.NewTrue())
	third := seocnd.Fork(context.NewFalse())

	// Assert
	assert.True(t, first.IsFeasible(), first.pc.String())
	assert.True(t, seocnd.IsFeasible(), seocnd.pc.String())
	assert.False(t, third.IsFeasible(), third.pc.String())
}

func TestBranchingPath(t *testing.T) {
	/*    a
	 *   / \
	 * x>0 x<10
	 *  |   |
	 *  b-<-c
	 *  |
	 */

	// Arrange
	context := z3.NewContext(z3.NewConfig())
	solver := context.NewSolver()
	x := context.NewConstant(z3.WithName("x"), context.IntegerSort())
	zero := context.NewInt(0, context.IntegerSort())
	ten := context.NewInt(10, context.IntegerSort())
	bc := z3.GT(x, zero)
	cc := z3.LT(x, ten)

	// Act
	a := NewGlobalPath(
		NewLexicalScope(context), context.NewTrue(),
	)
	b := a.Fork(bc)
	c := a.Fork(cc)

	b.Join(c)

	// Assert
	assert.True(t, a.IsFeasible(), a.pc.String())
	assert.True(t, c.IsFeasible(), c.pc.String())
	assert.True(t, b.IsFeasible(), b.pc.String())
	assert.True(t, solver.Proven(z3.Eq(a.pc, context.NewTrue())), a.pc.String())
	assert.True(t, solver.Proven(z3.Eq(b.pc, z3.Or(bc, cc))), c.pc.String())
	assert.True(t, solver.Proven(z3.Eq(c.pc, cc)), b.pc.String())
}

func TestBranchingValuations(t *testing.T) {
	/*    a
	 *   / \
	 * x>0 x<10
	 *  |   |
	 *  b-<-c
	 *  |
	 */

	// Arrange
	context := z3.NewContext(z3.NewConfig())
	solver := context.NewSolver()
	zero := context.NewInt(0, context.IntegerSort())
	ten := context.NewInt(10, context.IntegerSort())
	symbols := NewSymbolsMap(NewSymbolsFactory())
	xS, nS, mS := symbols.Insert("x"), symbols.Insert("n"), symbols.Insert("m")

	// Act
	a := NewGlobalPath(NewLexicalScope(context), context.NewTrue())
	x := a.scope.Define(xS, context.IntegerSort())
	a.scope.Declare(nS, zero)

	bc := z3.GT(x, zero)
	b := a.Fork(bc)
	b.scope.Declare(mS, ten)

	cc := z3.LT(x, ten)
	c := a.Fork(cc)
	c.scope.Assign(nS, ten)

	b.Join(c)

	// Assert
	assert.True(t, solver.Proven(z3.Eq(a.scope.Valuation(xS), zero)))
	assert.True(t, solver.Proven(z3.Eq(a.scope.Valuation(nS), zero)))
	assert.Nil(t, a.scope.Valuation(mS))

	assert.True(t, solver.Proven(z3.Eq(b.scope.Valuation(xS), zero)))
	assert.True(t, solver.Proven(z3.Eq(b.scope.Valuation(nS), z3.ITE(cc, ten, zero))))
	assert.True(t, solver.Proven(z3.Eq(b.scope.Valuation(mS), ten)))

	assert.True(t, solver.Proven(z3.Eq(c.scope.Valuation(xS), zero)))
	assert.True(t, solver.Proven(z3.Eq(c.scope.Valuation(nS), ten)))
	assert.Nil(t, c.scope.Valuation(mS))
}
