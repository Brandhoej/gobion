package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
*/
import "C"
import "github.com/Brandhoej/gobion/internal/z3"

type Numerals interface {
	valuation
	Int | Real
}

func Add[T Numerals](lhs T, rhs ...T) T {
	return NewValuation[T](z3.Add(lhs.ast(), asts[T](rhs...)...))
}

func Multiply[T Numerals](lhs T, rhs ...T) T {
	return NewValuation[T](z3.Multiply(lhs.ast(), asts[T](rhs...)...))
}

func Subtraction[T Numerals](lhs T, rhs ...T) T {
	return NewValuation[T](z3.Subtraction(lhs.ast(), asts[T](rhs...)...))
}

func Minus[T Numerals](operand T) T {
	return NewValuation[T](z3.Minus(operand.ast()))
}

func Division[T Numerals](lhs, rhs T) T {
	return NewValuation[T](z3.Division(lhs.ast(), rhs.ast()))
}

func Modulus[T Numerals](lhs, rhs T) T {
	return NewValuation[T](z3.Modulus(lhs.ast(), rhs.ast()))
}

func Remaninder[T Numerals](lhs, rhs T) T {
	return NewValuation[T](z3.Remaninder(lhs.ast(), rhs.ast()))
}

func Power[T Numerals](lhs, rhs T) T {
	return NewValuation[T](z3.Power(lhs.ast(), rhs.ast()))
}

func LT[T Numerals](lhs, rhs T) Bool {
	return newBool(z3.LT(lhs.ast(), rhs.ast()))
}

func LE[T Numerals](lhs, rhs T) Bool {
	return newBool(z3.LE(lhs.ast(), rhs.ast()))
}

func GT[T Numerals](lhs, rhs T) Bool {
	return newBool(z3.GT(lhs.ast(), rhs.ast()))
}

func GE[T Numerals](lhs, rhs T) Bool {
	return newBool(z3.GE(lhs.ast(), rhs.ast()))
}

func Divides[T Numerals](lhs, rhs T) Bool {
	return newBool(z3.Divides(lhs.ast(), rhs.ast()))
}

func IsInt(operand Real) Bool {
	return newBool(z3.IsInt(operand.ast()))
}
