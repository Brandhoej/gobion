package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

import "github.com/Brandhoej/gobion/internal/z3"

type Real struct {
	_ast  *z3.AST
	_sort *z3.Sort
}

func newReal(ast *z3.AST) Real {
	return Real{
		_ast:  ast,
		_sort: ast.Context().RealSort(),
	}
}

func (real Real) ast() *z3.AST {
	return real._ast
}

func (real Real) sort() *z3.Sort {
	return real._sort
}

func (real Real) String() string {
	return real._ast.String()
}
