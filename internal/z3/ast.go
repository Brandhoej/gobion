package z3

/*
#cgo CFLAGS: -I../../modules/z3
#cgo LDFLAGS: -L../../modules/z3 -lz3
#include "../../modules/z3/src/api/z3.h"
*/
import "C"
import "runtime"

// ASTKind is a general category of ASTs, such as numerals, applications, or sorts.
type ASTKind int

// The different kinds of Z3 AST (abstract syntax trees). That is, terms, formulas and types.
const (
	ASTKindApp        = ASTKind(C.Z3_APP_AST)        // Constant and applications
	ASTKindNumeral    = ASTKind(C.Z3_NUMERAL_AST)    // Numeral constants (excluding real algebraic numbers)
	ASTKindVar        = ASTKind(C.Z3_VAR_AST)        // Bound variables
	ASTKindQuantifier = ASTKind(C.Z3_QUANTIFIER_AST) // Quantifiers
	ASTKindSort       = ASTKind(C.Z3_SORT_AST)       // Sorts
	ASTKindFuncDecl   = ASTKind(C.Z3_FUNC_DECL_AST)  // Function declarations
	ASTKindUnknown    = ASTKind(C.Z3_UNKNOWN_AST)    // Z3 internal
)

// Abstract syntax tree node. That is, the data-structure used in Z3 to represent terms, formulas and types.
type AST struct {
	context *Context
	z3AST   C.Z3_ast
}

func (context *Context) wrapAST(z3AST C.Z3_ast) *AST {
	ast := &AST{
		context: context,
		z3AST:   z3AST,
	}

	// We force our own reference counting of the AST by using the specific rc function to create the context.
	C.Z3_inc_ref(context.z3Context, z3AST)
	runtime.SetFinalizer(ast, func(ast *AST) {
		// Make derement of reference counter atomic by wrapping it in a locked state.
		context.do(func() {
			C.Z3_dec_ref(context.z3Context, ast.z3AST)
		}, ast)
	})

	return ast
}

func (ast *AST) Context() *Context {
	return ast.context
}

func (ast *AST) Substitute(from, to []*AST) *AST {
	if len(from) != len(to) {
		panic("Substitution to/from must have the same length")
	}
	length := len(from)
	context := ast.context

	cFrom := make([]C.Z3_ast, length)
	cTo := make([]C.Z3_ast, length)

	for idx := range from {
		cFrom[idx] = from[idx].z3AST
		cTo[idx] = to[idx].z3AST
	}

	return compute(context, func() *AST {
		return context.wrapAST(
			C.Z3_substitute(
				context.z3Context,
				ast.z3AST, C.uint(length),
				&cFrom[0], &cTo[0],
			),
		)
	}, ast, from, to)
}

func (ast *AST) SubstituteVariables(to []*AST) *AST {
	length := len(to)
	context := ast.context
	cTo := make([]C.Z3_ast, 0, length)

	for idx := range to {
		cTo[idx] = to[idx].z3AST
	}

	return compute(context, func() *AST {
		return context.wrapAST(
			C.Z3_substitute_vars(
				context.z3Context, ast.z3AST,
				C.uint(length), &cTo[0],
			),
		)
	}, ast, to)
}

func (ast *AST) Simplify() *AST {
	return compute(ast.context, func() *AST {
		return ast.context.wrapAST(
			C.Z3_simplify(ast.context.z3Context, ast.z3AST),
		)
	}, ast)
}

// Convert the given AST node into a string.
//
// The result buffer is statically allocated by Z3. It will
// be automatically deallocated when Z3_del_context is invoked.
// So, the buffer is invalidated in the next call to Z3_ast_to_string.
func (ast *AST) String() string {
	return compute(ast.context, func() string {
		return C.GoString(C.Z3_ast_to_string(ast.context.z3Context, ast.z3AST))
	}, ast)
}

// Compares two ASTs (terms) for equality).
func (ast *AST) Equals(other *AST) bool {
	return compute(ast.context, func() bool {
		return bool(C.Z3_is_eq_ast(ast.context.z3Context, ast.z3AST, other.z3AST))
	}, ast, other)
}

// Return a hash code for the given AST.
// The hash code is structural but two different AST objects can map to the same hash.
// The result of Z3_get_ast_id returns an identifier that is unique over the
// set of live AST objects.
func (ast *AST) Hash() uint64 {
	return compute(ast.context, func() uint64 {
		return uint64(C.Z3_get_ast_hash(ast.context.z3Context, ast.z3AST))
	}, ast)
}

func (ast *AST) Sort() *Sort {
	return compute(ast.context, func() *Sort {
		return ast.context.wrapSort(
			C.Z3_get_sort(ast.context.z3Context, ast.z3AST),
		)
	}, ast)
}
