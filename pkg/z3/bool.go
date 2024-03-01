package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

import "github.com/Brandhoej/gobion/internal/z3"

// Represents an AST of the Bool sort.
type Bool struct {
	_ast  *z3.AST
	_sort *z3.Sort
}

func newBool(ast *z3.AST) Bool {
	return Bool{
		_ast:  ast,
		_sort: ast.Context().BooleanSort(),
	}
}

func (boolean Bool) ast() *z3.AST {
	return boolean._ast
}

func (boolean Bool) sort() *z3.Sort {
	return boolean._sort
}

func (boolean Bool) String() string {
	return boolean._ast.String()
}
