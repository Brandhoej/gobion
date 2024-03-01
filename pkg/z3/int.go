package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

import "github.com/Brandhoej/gobion/internal/z3"

type Int struct {
	_ast  *z3.AST
	_sort *z3.Sort
}

func newInt(ast *z3.AST) Int {
	return Int{
		_ast:  ast,
		_sort: ast.Context().IntegerSort(),
	}
}

func (integer Int) ast() *z3.AST {
	return integer._ast
}

func (integer Int) sort() *z3.Sort {
	return integer._sort
}

func (integer Int) String() string {
	return integer._ast.String()
}
