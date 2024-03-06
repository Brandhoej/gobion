package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Scope interface {
	// Returns true if the identifier is declared in the scope and not one of it ancestors.
	IsLocal(identifier string) bool

	// Returns true if the scope is the global scope of the program (i.e., it does not have a parent scope).
	IsGlobal() bool

	// Declares a variable with a zero-valuation.
	Define(identifier string, sort *z3.Sort) (variable *z3.AST)

	// Declares a variable with the valuation.
	Declare(identifier string, valuation *z3.AST) (variable *z3.AST)

	// Assigned a new valuation to a varaible.
	// This updates the variable stored for the identifier to ensure Single Static Assignment (SSA).
	Assign(identifier string, valuation *z3.AST) (variable *z3.AST)

	// Looks up a variable and returns it.
	Variable(identifier string) (variable *z3.AST, exists bool)

	// Looks up a valuation. If the variable is bound that the valuation is the bound variable itself.
	// However, if the variable is not bound then the valuation is the reduced variables.
	Valuation(identifier string) (valuation *z3.AST, exists bool)

	// Returns the set of all local identifiers of the scope.
	Identifiers() []string
}
