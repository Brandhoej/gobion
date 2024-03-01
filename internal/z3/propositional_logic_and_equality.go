package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

func Not(operand *AST) *AST {
	return unary(
		func(context C.Z3_context, operand C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_not(context, operand)
		}, operand,
	)
}

func And(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_and(context, length, &operands[0])
		}, lhs, rhs...,
	)
}

func Or(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_or(context, length, &operands[0])
		}, lhs, rhs...,
	)
}

func Xor(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_xor(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func IFF(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_iff(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func ITE(condition, consequence, alternative *AST) *AST {
	return ternary(
		func(context C.Z3_context, a, b, c C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_ite(context, a, b, c)
		}, condition, consequence, alternative,
	)
}

func Implies(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_implies(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Eq(lhs, rhs *AST) *AST {
	return binary(
		func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_eq(context, lhs, rhs)
		}, lhs, rhs,
	)
}

func Distinct(lhs *AST, rhs ...*AST) *AST {
	return nary(
		func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast {
			return C.Z3_mk_distinct(context, length, &operands[0])
		}, lhs, rhs...,
	)
}
