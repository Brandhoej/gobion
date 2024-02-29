package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"sync"
)

// Manager of all other Z3 objects, global configuration options, etc.
type Context struct {
	// Manager of all other Z3 objects, global configuration options, etc.
	z3Context C.Z3_context

	// mutex protects AST reference counts and the context's last
	// error. Use Context.do to acquire this around a Z3 operation
	// and panic if the operation has an error status.
	//
	// This is necessary as we used Z3_mk_context_rc with our own reference counting.
	//   This reference counting is important when performing AST operations.
	mutex sync.Mutex
}

// Create a context using the given configuration.
func NewContext(config *Config) *Context {
	context := &Context{
		// This function is similar to Z3_mk_context. However,
		// in the context returned by this function, the user
		// is responsible for managing Z3_ast reference counters.
		// Managing reference counters is a burden and error-prone,
		// but allows the user to use the memory more efficiently.
		// The user must invoke Z3_inc_ref for any Z3_ast returned
		// by Z3, and Z3_dec_ref whenever the Z3_ast is not needed
		// anymore. This idiom is similar to the one used in
		// BDD (binary decision diagrams) packages such as CUDD.
		z3Context: C.Z3_mk_context_rc(config.z3Config),
	}

	// Before GC of the context we want to delete the C unmanaged context object.
	runtime.SetFinalizer(context, func(context *Context) {
		C.Z3_del_context(context.z3Context)
	})

	return context
}

// Interrupt the execution of a Z3 procedure.
// This procedure can be used to interrupt: solvers, simplifiers and tactics.
func (context *Context) Interrupt() {
	C.Z3_interrupt(context.z3Context)

	// It might be the intention to interrupt the context and then stop the program.
	// In these cases the references to the context would be removed and the finaliser would be called.
	// The finaliser of the context would delete the z3Context. In some cases the context would be finalised
	// before the interrupt is called or the finalisation happens concurrently to the finalisation.
	// This can lead to non-deterministic and wrong behvaiour. Therefore, in order to await the finalisation
	// of the context we keep the object alive.
	runtime.KeepAlive(context)
}

// Aquires the mutex lock necessary for performing AST operations from the context.
func (context *Context) do(action func(), keeps ...any) {
	context.mutex.Lock()
	defer func() {
		context.mutex.Unlock()
		for _, keep := range keeps {
			runtime.KeepAlive(keep)
		}
	}()
	action()
}

func compute[T any](context *Context, function func() T, keeps ...any) T {
	var value T

	context.do(func() {
		value = function()
	}, keeps...)

	return value
}

func (context *Context) NewConstant(symbolFactory SymbolFactory, sort *Sort) *AST {
	symbol := symbolFactory(context)
	return compute[*AST](context, func() *AST {
		return context.wrapAST(
			C.Z3_mk_const(context.z3Context, symbol.z3Symbol, sort.z3Sort),
		)
	}, sort, symbolFactory)
}

func (context *Context) NewReal(numerator, denominator int) *AST {
	return compute[*AST](context, func() *AST {
		return context.wrapAST(
			C.Z3_mk_real(context.z3Context, C.int(numerator), C.int(denominator)),
		)
	})
}

func (context *Context) NewInt(value int, sort *Sort) *AST {
	return compute[*AST](context, func() *AST {
		return context.wrapAST(
			C.Z3_mk_int(context.z3Context, C.int(value), sort.z3Sort),
		)
	}, sort)
}

func (context *Context) NewTrue() *AST {
	return context.wrapAST(
		C.Z3_mk_true(context.z3Context),
	)
}

func (context *Context) NewFalse() *AST {
	return context.wrapAST(
		C.Z3_mk_false(context.z3Context),
	)
}
