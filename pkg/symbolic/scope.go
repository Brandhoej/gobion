package symbolic

import "github.com/Brandhoej/gobion/internal/z3"

type Scope interface {
	IsLocal(identifier string) bool
	IsGlobal() bool
	
	Bind(index uint, sort *z3.Sort) (variable *z3.AST)
	Define(identifier string, sort *z3.Sort) (variable *z3.AST)
	Declare(identifier string, valuation *z3.AST) (variable *z3.AST)
	Assign(identifier string, valuation *z3.AST) (variable *z3.AST)

	Variable(identifier string) (variable *z3.AST, exists bool)
	Valuation(identifier string) (valuation *z3.AST, exists bool)
}