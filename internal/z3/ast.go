package z3

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
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
	noEq
}

func (context *Context) wrapAST(z3AST C.Z3_ast) *AST {
	ast := &AST{
		context: context,
		z3AST:   z3AST,

		// AST is not directly comparable and Equals should be used instead.
		noEq: noEq{},
	}

	// We force our own reference counting of the AST by using the specific rc function to create the context.
	C.Z3_inc_ref(context.z3Context, z3AST)
	runtime.SetFinalizer(ast, func(ast *AST) {
		// Make derement of reference counter atomic by wrapping it in a locked state.
		context.do(func() {
			C.Z3_dec_ref(context.z3Context, ast.z3AST)
		})
	})

	return ast
}

func (ast *AST) Context() *Context {
	return ast.context
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
