package symbolic

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
)

type GoPath struct {
	parent  *GoPath
	scope   *GoScope
	context *z3.Context
	pc      *z3.AST
}

func NewGoGlobalPath(context *z3.Context) *GoPath {
	return &GoPath{
		parent:  nil,
		scope:   NewGoGlobalScope(),
		context: context,
		pc:      context.NewTrue(),
	}
}

func (path *GoPath) Branch(condition *z3.AST) *GoPath {
	solver := path.context.NewSolver()

	pc := z3.And(path.pc, condition).Simplify()

	// Check if new pc is a contradiction if so we dont branch.
	solver.Push()
	solver.Assert(pc)
	for _, identifier := range path.scope.Identifiers() {
		variable, _ := path.scope.Variable(identifier)
		valuation, _ := path.scope.Valuation(identifier)
		solver.Assert(z3.Eq(variable, valuation))
	}
	hasSolution := solver.Check().IsTrue()
	solver.Pop(1)

	if !hasSolution {
		return nil
	}

	return &GoPath{
		parent:  path,
		scope:   path.scope.Child(),
		context: path.context,
		pc:      pc,
	}
}

// Merges an if-then-else branch.
func (consequence *GoPath) MergeITE(alternative *GoPath) *GoPath {
	parent := consequence.parent

	junction := &GoPath{
		parent:  parent,
		scope:   parent.scope,
		context: parent.context,
		pc:      parent.pc,
	}

	// We only merge the variables present in the parent path as the variables declared in the scopes
	// of either the consequence and alternative should not be visible after the corresponding blocks.
	// The variables that have different values in the branches need to be assigned with ITE.
	for _, identifier := range junction.scope.Identifiers() {
		consequenceValue, _ := consequence.scope.Valuation(identifier)
		alternativeValue, _ := alternative.scope.Valuation(identifier)

		if consequenceValue != alternativeValue {
			ite := z3.ITE(consequence.pc, consequenceValue, alternativeValue).Simplify()
			junction.scope.Assign(identifier, ite)
		}
	}

	return junction
}

// Meges an if-then branching.
func (consequence *GoPath) MergeIT() *GoPath {
	return consequence.MergeITE(consequence.parent)
}

func (path *GoPath) String() string {
	var builder strings.Builder
	scope := path.scope.String()
	builder.WriteString(scope)
	builder.WriteRune('\n')
	pc := path.pc.String()
	builder.WriteString(fmt.Sprintf("pc=%s", pc))
	return builder.String()
}
