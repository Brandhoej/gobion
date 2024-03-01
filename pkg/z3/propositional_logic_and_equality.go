package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"
import "github.com/Brandhoej/gobion/internal/z3"

func Not(operand Bool) Bool {
	return newBool(z3.Not(operand.ast()))
}

func And(lhs Bool, rhs ...Bool) Bool {
	return newBool(z3.And(lhs.ast(), asts[Bool](rhs...)...))
}

func Or(lhs Bool, rhs ...Bool) Bool {
	return newBool(z3.Or(lhs.ast(), asts[Bool](rhs...)...))
}

func Xor(lhs, rhs Bool) Bool {
	return newBool(z3.Xor(lhs.ast(), rhs.ast()))
}

func IFF(lhs, rhs Bool) Bool {
	return newBool(z3.IFF(lhs.ast(), rhs.ast()))
}

func ITE[T valuation](condition Bool, consequence, alternative T) T {
	return NewValuation[T](z3.ITE(condition.ast(), consequence.ast(), alternative.ast()))
}

func Implies(lhs, rhs Bool) Bool {
	return newBool(z3.Implies(lhs.ast(), rhs.ast()))
}

func Eq[T valuation](lhs, rhs T) Bool {
	return newBool(z3.Eq(lhs.ast(), rhs.ast()))
}

func Distinct[T valuation](lhs T, rhs ...T) T {
	return NewValuation[T](z3.Distinct(lhs.ast(), asts[T](rhs...)...))
}
