package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

func Add(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_add(context, length, &operands[0])
		}, lhs, rhs...,
	)
}

func Multiply(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_mul(context, length, &operands[0])
		}, lhs, rhs...,
	)
}

func Subtract(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_sub(context, length, &operands[0])
		}, lhs, rhs...,
	)
}

func Minus(operand *AST) *AST {
	return unary(
		func(context C.Z3_context, operand C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_unary_minus(context, operand)
		}, operand,
	)
}

func Divide(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_div(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Modulus(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_mod(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Remaninder(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_rem(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Power(base, exponent *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_power(context, lhs, rhs)
		}, base, exponent,
	)
}

func LT(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_lt(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func LE(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_le(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func GT(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_gt(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func GE(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_ge(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Divides(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_divides(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func IsInt(operand *AST) *AST {
	return unary(
		func(context C.Z3_context, operand C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_is_int(context, operand)
		}, operand,
	)
}
