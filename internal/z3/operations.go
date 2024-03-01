package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

func unary(
	operation func(context C.Z3_context, operand C.Z3_ast) C.Z3_ast, operand *AST,
) *AST {
	return compute[*AST](operand.context, func() *AST {
		return operand.context.wrapAST(
			operation(operand.context.z3Context, operand.z3AST),
		)
	})
}

func binary(
	operation func(context C.Z3_context, lhs, rhs C.Z3_ast) C.Z3_ast,
	lhs, rhs *AST,
) *AST {
	return compute[*AST](lhs.context, func() *AST {
		return lhs.context.wrapAST(
			operation(
				lhs.context.z3Context,
				lhs.z3AST,
				rhs.z3AST,
			),
		)
	})
}

func ternary(
	operation func(context C.Z3_context, a, b, c C.Z3_ast) C.Z3_ast,
	a, b, c *AST,
) *AST {
	return compute[*AST](a.context, func() *AST {
		return a.context.wrapAST(
			operation(a.context.z3Context,
				a.z3AST,
				b.z3AST,
				c.z3AST,
			),
		)
	})
}

func nary(
	operation func(context C.Z3_context, length C.uint, operands ...C.Z3_ast) C.Z3_ast,
	operand *AST, operands ...*AST,
) *AST {
	// Create the n-ary operand array.
	args := make([]C.Z3_ast, len(operands)+1)
	args[0] = operand.z3AST
	for i, operand := range operands {
		args[i+1] = operand.z3AST
	}

	return compute[*AST](operand.context, func() *AST {
		return operand.context.wrapAST(
			operation(
				operand.context.z3Context,
				C.uint(len(args)),
				args...,
			),
		)
	})
}
