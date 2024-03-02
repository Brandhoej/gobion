package symbolic

import (
	"fmt"
	"strings"

	"github.com/Brandhoej/gobion/internal/z3"
)

type Path interface {
	Declare(identifier string, value *z3.AST)
	Assign(lhs, rhs *z3.AST)
	Branch(condition *z3.AST)
}

type GoPath struct {
	parent  *GoPath
	scope   *GoScope[*z3.AST]
	context *z3.Context
	pc      *z3.AST
}

func NewGoGlobalPath(context *z3.Context) *GoPath {
	return &GoPath{
		parent: nil,
		scope: NewGoGlobalScope[*z3.AST](),
		context: context,
		pc: context.NewTrue(),
	}
}

func (path *GoPath) Branch(condition *z3.AST) *GoPath {
	solver := path.context.NewSolver()

	pc := z3.And(path.pc, condition).Simplify()

	// Check if new pc is a contradiction if so we dont branch.
	solver.Push()
	solver.Assert(pc)
	for _, symbol := range path.scope.Symbols() {
		variable, _ := path.scope.Variable(symbol)
		valuation, _ := path.scope.Valuation(symbol)
		solver.Assert(z3.Eq(variable, valuation))
	}
	hasSolution := solver.Check().IsTrue()
	solver.Pop(1)

	if !hasSolution {
		return nil
	}

	return &GoPath{
		parent: path,
		scope:   path.scope.Child(),
		context: path.context,
		pc:      pc,
	}
}

// Merges an if-then-else branch.
func (consequence *GoPath) MergeITE(alternative *GoPath) *GoPath {
	parent := consequence.parent

	junction := &GoPath{
		parent: parent,
		scope:   parent.scope,
		context: parent.context,
		pc:      parent.pc,
	}

	// We only merge the variables present in the parent path as the variables declared in the scopes
	// of either the consequence and alternative should not be visible after the corresponding blocks.
	// The variables that have different values in the branches need to be assigned with ITE.
	for _, symbol := range junction.scope.Symbols() {
		consequenceValue, _ := consequence.scope.Valuation(symbol)
		alternativeValue, _ := alternative.scope.Valuation(symbol)

		if consequenceValue != alternativeValue {
			ite := z3.ITE(consequence.pc, consequenceValue, alternativeValue).Simplify()
			junction.scope.Assign(symbol, ite)
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
	builder.WriteString(path.scope.String())
	builder.WriteRune('\n')
	builder.WriteString(fmt.Sprintf("pc=%s", path.pc.String()))
	return builder.String()
}
