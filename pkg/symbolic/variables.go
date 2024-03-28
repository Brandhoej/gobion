package symbolic

import (
	"fmt"

	"github.com/Brandhoej/gobion/internal/z3"
)

type Variables interface {
	Declare(symbol Symbol, identifier string, sort *z3.Sort) (variable *z3.AST)
	Advance(symbol Symbol) (variable *z3.AST, exists bool)
	Variable(symbol Symbol) (variable *z3.AST, exists bool)
}

type variable struct {
	name string
	identifier string
	counter int
	ast *z3.AST
}

func newVariable(name string, sort *z3.Sort) *variable {
	identifier := fmt.Sprintf("%s_%v", name, 0)
	return &variable{
		name: name,
		identifier: identifier,
		counter: 0,
		ast: sort.Context().NewConstant(
			z3.WithName(identifier), sort,
		),
	}
}

func (variable *variable) advance() {
	variable.counter += 1
	variable.identifier = fmt.Sprintf("%s_%v", variable.name, variable.counter)
	variable.ast = variable.ast.Context().NewConstant(
		z3.WithName(variable.identifier), variable.ast.Sort(),
	)
}

type VariablesMap struct {
	variables   map[Symbol]*variable
}

func NewVariablesMap() *VariablesMap {
	return &VariablesMap{
		variables: map[Symbol]*variable{},
	}
}

func (mapping *VariablesMap) Declare(symbol Symbol, name string, sort *z3.Sort) *z3.AST {
	if _, exists := mapping.variables[symbol]; exists {
		return nil
	}

	variable := newVariable(name, sort)
	mapping.variables[symbol] = variable

	return variable.ast
}

func (mapping *VariablesMap) Advance(symbol Symbol) (*z3.AST, bool) {
	if variable, exists := mapping.variables[symbol]; exists {
		variable.advance()
		return variable.ast, false
	}
	return nil, false
}

func (mapping *VariablesMap) Identifier(symbol Symbol) (string, bool) {
	if variable, exists := mapping.variables[symbol]; exists {
		return variable.identifier, exists
	}
	return "", false
}

func (mapping *VariablesMap) NameOf(symbol Symbol) (string, bool) {
	if variable, exists := mapping.variables[symbol]; exists {
		return variable.name, exists
	}
	return "", false
}

func (mapping *VariablesMap) Variable(symbol Symbol) (*z3.AST, bool) {
	if variable, exists := mapping.variables[symbol]; exists {
		return variable.ast, exists
	}
	return nil, false
}
