package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

// Kind of AST used to represent function symbols.
type FunctionDeclaration struct {
	context               *Context
	z3FunctionDeclaration C.Z3_func_decl
}

func (context *Context) wrapFunctionDeclaration(function C.Z3_func_decl) *FunctionDeclaration {
	return &FunctionDeclaration{
		context:               context,
		z3FunctionDeclaration: function,
	}
}

func (function *FunctionDeclaration) Application(arguments []*AST) *AST {
	return compute(function.context, func() *AST {
		args := make([]C.Z3_ast, len(arguments))
		for i, operand := range arguments {
			args[i] = operand.z3AST
		}

		return function.context.wrapAST(
			C.Z3_mk_app(
				function.context.z3Context,
				function.z3FunctionDeclaration,
				C.uint(len(arguments)),
				&args[0],
			),
		)
	}, function, arguments)
}

func (function *FunctionDeclaration) AST() *AST {
	return compute(function.context, func() *AST {
		return function.context.wrapAST(
			C.Z3_func_decl_to_ast(
				function.context.z3Context,
				function.z3FunctionDeclaration,
			),
		)
	}, function)
}
