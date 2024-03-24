package golang

import (
	"go/ast"

	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
	"github.com/Brandhoej/gobion/pkg/scfg/scfg"
)

type factory struct {
	flow *cfg.Graph[ast.Stmt, ast.Expr]
	scopes *scfg.Graph[ast.Stmt, ast.Expr]
	block int
}

func new() *factory {
	flow := cfg.New[ast.Stmt, ast.Expr]()
	scopes := scfg.New(flow)
	scopes.Into(scopes.Global(), flow.Entry())
	return &factory{
		flow: flow,
		scopes: scopes,
		block: flow.Entry(),
	}
}

func (factory *factory) GetBlock() int {
	return factory.block
}

func (factory *factory) GetCondition() int {
	_, condition := factory.flow.Block(factory.block)
	return condition
}

func (factory *factory) SetBlock(block int) int {
	factory.block = block
	return factory.block
}

func (factory *factory) Append(statements ...ast.Stmt) {
	factory.flow.Append(factory.block, statements...)
}

func (factory *factory) Proceed(scope int) int {
	block, _ := factory.flow.NewBlock()
	factory.SetBlock(block)
	factory.scopes.Into(scope, block)
	return block
}

func (factory *factory) Condition(block int, expression ast.Expr) int {
	condition := factory.flow.NewCondition(expression)
	factory.flow.JumpTo(block, condition)
	return condition
}

func (factory *factory) JumpTo(block, condition int) {
	factory.flow.JumpTo(block, condition)
}

func (factory *factory) JumpFrom(condition, jump int) {
	factory.flow.JumpFrom(condition, jump)
}

func (factory *factory) Connect(source, destination int, expression ast.Expr) int {
	jump := factory.flow.NewConditionalJump(expression, destination)
	_, condition := factory.flow.Block(source)
	factory.flow.JumpFrom(condition, jump)
	factory.SetBlock(destination)
	return destination
}

func (factory *factory) ConnectFrom(source int, expression ast.Expr) int {
	return factory.Connect(source, factory.GetBlock(), expression)
}

func (factory *factory) ConnectTo(destination int, expression ast.Expr) int {
	return factory.Connect(factory.GetBlock(), destination, expression)
}

func (factory *factory) ZoomIn(source int, expression ast.Expr) int {
	scope, exists := factory.scopes.ScopeWith(source)
	if !exists {
		panic("Cannot zoom in from a block that is not in a scope")
	}
	inner := factory.scopes.ZoomIn(scope)
	proceeding := factory.Proceed(inner)
	factory.Connect(source, proceeding, expression)
	return inner
}

func (factory *factory) ZoomOut(block int) int {
	scope, exists := factory.scopes.ScopeWith(block)
	if !exists {
		panic("Cannot zoom out from a block that is not in a scope")
	}

	outer, _ := factory.scopes.ZoomOut(scope)
	return outer
}

func (factory *factory) statement(statement ast.Stmt) {
	if statement == nil || factory.block == factory.flow.Exit() {
		return
	}

	scope, exists := factory.scopes.ScopeWith(factory.GetBlock())
	if !exists {
		panic("Mismatch between current block scope and scope")
	}

	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		for _, statement := range cast.List {
			factory.statement(statement)
		}
	case *ast.IfStmt:
		if cast.Init != nil {
			factory.ZoomIn(factory.GetBlock(), nil)
			factory.statement(cast.Init)
		}
		
		front := factory.GetBlock()
		factory.Condition(front, cast.Cond)

		hasConsequence := len(cast.Body.List) > 0
		hasAlternative := cast.Else != nil && len(cast.Else.(*ast.BlockStmt).List) > 0

		if hasConsequence {
			factory.ZoomIn(front, ast.NewIdent("true"))
			factory.statement(cast.Body)
		}
		consequence := factory.GetBlock()

		if hasAlternative {
			factory.ZoomIn(front, ast.NewIdent("false"))
			factory.statement(cast.Else)
		}
		alternative := factory.GetBlock()

		proceeding := factory.Proceed(scope)
		if hasConsequence && factory.flow.Exit() != consequence {
			factory.Connect(consequence, proceeding, nil)
		}
		if hasAlternative && factory.flow.Exit() != alternative {
			factory.Connect(alternative, proceeding, nil)
		} else {
			factory.Connect(front, proceeding, ast.NewIdent("false"))
		}
	case *ast.ForStmt:
		outer := factory.ZoomIn(factory.GetBlock(), nil)

		if cast.Init != nil {
			factory.statement(cast.Init)
		}

		condition := cast.Cond
		if condition == nil {
			condition = ast.NewIdent("true")
		}

		entryCondition := factory.Condition(factory.GetBlock(), condition)
	
		factory.ZoomIn(factory.GetBlock(), ast.NewIdent("true"))
		factory.statement(cast.Body)
		exit := factory.GetBlock()

		if cast.Post != nil {
			factory.Proceed(outer)
			factory.ConnectFrom(exit, nil)
			factory.statement(cast.Post)
		}
		
		factory.JumpTo(factory.GetBlock(), entryCondition)

		proceeding := factory.Proceed(scope)
		factory.JumpFrom(
			entryCondition, 
			factory.flow.NewConditionalJump(ast.NewIdent("false"), proceeding),
		)
	case *ast.ReturnStmt:
		factory.Append(statement)
		factory.Connect(factory.block, factory.flow.Exit(), nil)
	default:
		factory.Append(statement)
	}
}

func SCFG(function *ast.FuncDecl) *scfg.Graph[ast.Stmt, ast.Expr] {
	factory := new()
	factory.statement(function.Body)
	return factory.scopes
}