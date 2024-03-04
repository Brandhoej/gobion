package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

type ASTVector struct {
	context *Context
	vector  C.Z3_ast_vector
}

func (context *Context) wrapASTVector(vector C.Z3_ast_vector) *ASTVector {
	return &ASTVector{
		context: context,
		vector:  vector,
	}
}

func (asts *ASTVector) Length() uint {
	return uint(C.Z3_ast_vector_size(asts.context.z3Context, asts.vector))
}

func (asts *ASTVector) Get(index uint) *AST {
	return asts.context.wrapAST(
		C.Z3_ast_vector_get(
			asts.context.z3Context,
			asts.vector, C.uint(index),
		),
	)
}
