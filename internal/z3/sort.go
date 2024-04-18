package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"

// Kind of AST used to represent types.
type Sort struct {
	context *Context
	z3Sort  C.Z3_sort
}

func (sort *Sort) AST() *AST {
	return sort.context.wrapAST(
		C.Z3_sort_to_ast(sort.context.z3Context, sort.z3Sort),
	)
}

func (sort *Sort) Zero() (zero *AST) {
	switch sort.Kind() {
	case KindInt:
		zero = sort.context.NewInt(0, sort)
	case KindBoolean:
		zero = sort.context.NewFalse()
	}
	return zero
}

func (sort *Sort) Context() *Context {
	return sort.context
}

func (sort *Sort) Kind() Kind {
	return Kind(C.Z3_get_sort_kind(sort.context.z3Context, sort.z3Sort))
}

func (sort *Sort) SameAs(others ...*Sort) bool {
	for _, other := range others {
		if *sort != *other {
			return false
		}
	}
	return true
}

func (context *Context) wrapSort(z3Sort C.Z3_sort) *Sort {
	return &Sort{
		context: context,
		z3Sort:  z3Sort,
	}
}

func (context *Context) BooleanSort() *Sort {
	return context.wrapSort(
		C.Z3_mk_bool_sort(context.z3Context),
	)
}

func (context *Context) IntegerSort() *Sort {
	return context.wrapSort(
		C.Z3_mk_int_sort(context.z3Context),
	)
}

func (context *Context) RealSort() *Sort {
	return context.wrapSort(
		C.Z3_mk_real_sort(context.z3Context),
	)
}
