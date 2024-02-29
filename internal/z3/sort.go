package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
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
