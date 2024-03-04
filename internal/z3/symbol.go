package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Symbol struct {
	context  *Context
	z3Symbol C.Z3_symbol
}

type SymbolFactory func(context *Context) Symbol

func WithName(name string) SymbolFactory {
	return func(context *Context) Symbol {
		return context.NewStringSymbol(name)
	}
}

func WithInt(value int) SymbolFactory {
	return func(context *Context) Symbol {
		return context.NewIntegerSymbol(value)
	}
}

func WithSymbol(symbol Symbol) SymbolFactory {
	return func(context *Context) Symbol {
		return symbol
	}
}

func (context *Context) NewStringSymbol(identifier string) Symbol {
	// Allocate an unmanged string and make sure it is freed.
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))

	return compute(context, func() Symbol {
		return Symbol{
			context:  context,
			z3Symbol: C.Z3_mk_string_symbol(context.z3Context, cIdentifier),
		}
	})
}

func (context *Context) NewIntegerSymbol(value int) Symbol {
	return compute(context, func() Symbol {
		return Symbol{
			context:  context,
			z3Symbol: C.Z3_mk_int_symbol(context.z3Context, C.int(value)),
		}
	})
}
