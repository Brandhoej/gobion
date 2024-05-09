package golang

import (
	"go/ast"
	"go/token"

	"github.com/Brandhoej/gobion/pkg/scfg/cfg"
	"github.com/Brandhoej/gobion/pkg/scfg/scfg"
	"github.com/Brandhoej/gobion/pkg/structures"
)

type factory struct {
	flow         *cfg.Graph[ast.Stmt, ast.Expr]
	scopes       *scfg.Graph[ast.Stmt, ast.Expr]
	labels       map[string]int
	gotos        map[string][]int
	terminated   map[int]any
	loopEnds     structures.Stack[int]
	loopExits    []int
	fallthroughs []int
	block        int
	withScopes   bool
}

func new() *factory {
	flow := cfg.New[ast.Stmt, ast.Expr]()
	scopes := scfg.New(flow)
	scopes.Into(scopes.Global(), flow.Entry())
	return &factory{
		flow:         flow,
		scopes:       scopes,
		block:        flow.Entry(),
		labels:       map[string]int{},
		gotos:        map[string][]int{},
		terminated:   map[int]any{},
		loopEnds:     make(structures.Stack[int], 0),
		loopExits:    make(structures.Stack[int], 0),
		fallthroughs: make([]int, 0),
	}
}

func (factory *factory) getBlock() int {
	return factory.block
}

func (factory *factory) setBlock(block int) int {
	factory.block = block
	return factory.block
}

func (factory *factory) append(statements ...ast.Stmt) {
	factory.flow.Append(factory.block, statements...)
}

func (factory *factory) proceed(scope int) int {
	block, _ := factory.flow.NewBlock()
	factory.setBlock(block)
	factory.scopes.Into(scope, block)
	return block
}

func (factory *factory) condition(block int, expression ast.Expr) int {
	condition := factory.flow.NewCondition(expression)
	factory.flow.JumpTo(block, condition)
	return condition
}

func (factory *factory) jumpTo(block, condition int) {
	factory.flow.JumpTo(block, condition)
}

func (factory *factory) jumpFrom(condition, jump int) {
	factory.flow.JumpFrom(condition, jump)
}

func (factory *factory) connect(source, destination int, expression ast.Expr) int {
	if source == factory.flow.Exit() {
		return source
	}

	jump := factory.flow.NewConditionalJump(expression, destination)
	_, condition := factory.flow.Block(source)
	factory.flow.JumpFrom(condition, jump)
	factory.setBlock(destination)
	return destination
}

func (factory *factory) connectFrom(source int, expression ast.Expr) int {
	return factory.connect(source, factory.getBlock(), expression)
}

func (factory *factory) zoomIn(source int, expressions ...ast.Expr) int {
	scope, exists := factory.scopes.ScopeWith(source)
	if !exists {
		panic("Cannot zoom in from a block that is not in a scope")
	}

	if !factory.withScopes {
		hasExpression := false
		for idx := range expressions {
			if expressions[idx] != nil {
				hasExpression = true
				break
			}
		}

		if hasExpression {
			proceeding := factory.proceed(scope)
			for idx := range expressions {
				factory.connect(source, proceeding, expressions[idx])
			}
		}

		return scope
	}

	inner := factory.scopes.ZoomIn(scope)
	proceeding := factory.proceed(inner)
	for idx := range expressions {
		factory.connect(source, proceeding, expressions[idx])
	}
	return inner
}

func (factory *factory) enterLoop(block int) {
	factory.loopEnds.Push(block)
}

func (factory *factory) exitLoop(proceeding int) {
	factory.loopEnds.Pop()

	for _, block := range factory.loopExits {
		factory.connect(block, proceeding, nil)
	}
	factory.loopExits = []int{}
}

func (factory *factory) isTerminated(block int) bool {
	_, termianted := factory.terminated[block]
	return termianted || factory.flow.Exit() == block
}

func (factory *factory) statements(statements []ast.Stmt) {
	for idx := range statements {
		factory.statement(statements[idx])
	}
}

func (factory *factory) statement(statement ast.Stmt) {
	if statement == nil || factory.block == factory.flow.Exit() {
		return
	}

	_, terminated := factory.terminated[factory.getBlock()]

	scope, exists := factory.scopes.ScopeWith(factory.getBlock())
	if !exists {
		panic("Mismatch between current block scope and scope")
	}

	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		if terminated {
			return
		}

		for _, statement := range cast.List {
			factory.statement(statement)
		}
	case *ast.IfStmt:
		if terminated {
			return
		}

		if cast.Init != nil {
			factory.zoomIn(factory.getBlock(), nil)
			factory.statement(cast.Init)
		}

		front := factory.getBlock()
		factory.condition(front, cast.Cond)

		hasConsequence := len(cast.Body.List) > 0
		hasAlternative := cast.Else != nil && len(cast.Else.(*ast.BlockStmt).List) > 0

		if hasConsequence {
			factory.zoomIn(front, ast.NewIdent("true"))
			factory.statement(cast.Body)
		}
		consequence := factory.getBlock()

		if hasAlternative {
			factory.zoomIn(front, ast.NewIdent("false"))
			factory.statement(cast.Else)
		}
		alternative := factory.getBlock()

		proceeding := factory.proceed(scope)
		if hasConsequence && !factory.isTerminated(consequence) {
			factory.connect(consequence, proceeding, nil)
		}
		if hasAlternative && !factory.isTerminated(alternative) {
			factory.connect(alternative, proceeding, nil)
		} else {
			factory.connect(front, proceeding, ast.NewIdent("false"))
		}
	case *ast.SwitchStmt:
		if terminated {
			return
		}

		if cast.Init != nil {
			factory.zoomIn(factory.getBlock(), nil)
			factory.statement(cast.Init)
		}

		front := factory.getBlock()
		factory.condition(front, cast.Tag)

		ends := make([]int, 0, len(cast.Body.List))
		for Cidx, clause := range cast.Body.List {
			clause := clause.(*ast.CaseClause)
			factory.zoomIn(front, clause.List...)
			factory.statements(clause.Body)
			end := factory.getBlock()
			ends = append(ends, end)

			if Cidx > 0 {
				for idx := range factory.fallthroughs {
					factory.connect(factory.fallthroughs[idx], end, nil)
				}
				if len(factory.fallthroughs) > 0 {
					factory.terminated[ends[Cidx-1]] = nil
				}
				factory.fallthroughs = []int{}
			}
		}

		allEndsTerminated := true
		for idx := range ends {
			if termianted := factory.isTerminated(ends[idx]); !termianted {
				allEndsTerminated = false
				break
			}
		}

		if !allEndsTerminated {
			proceeding := factory.proceed(scope)
			for _, end := range ends {
				if terminated := factory.isTerminated(end); !terminated {
					factory.connect(end, proceeding, nil)
				}
			}
		}
	case *ast.ForStmt:
		if terminated {
			return
		}

		outer := factory.zoomIn(factory.getBlock(), nil)

		if cast.Init != nil {
			factory.statement(cast.Init)
		}

		condition := cast.Cond
		if condition == nil {
			condition = ast.NewIdent("true")
		}

		entryCondition := factory.condition(factory.getBlock(), condition)

		factory.zoomIn(factory.getBlock(), ast.NewIdent("true"))
		factory.statement(cast.Body)
		exit := factory.getBlock()

		if cast.Post != nil {
			factory.proceed(outer)
			factory.connectFrom(exit, nil)
			factory.statement(cast.Post)
		}

		factory.enterLoop(exit)
		factory.jumpTo(factory.getBlock(), entryCondition)

		proceeding := factory.proceed(scope)
		factory.exitLoop(proceeding)
		factory.jumpFrom(
			entryCondition,
			factory.flow.NewConditionalJump(ast.NewIdent("false"), proceeding),
		)
	case *ast.ReturnStmt:
		if terminated {
			return
		}

		factory.append(statement)
		factory.connect(factory.getBlock(), factory.flow.Exit(), nil)
	case *ast.BranchStmt:
		if terminated {
			return
		}

		block := factory.getBlock()
		factory.terminated[block] = nil

		switch cast.Tok {
		case token.BREAK:
			factory.loopExits = append(factory.loopExits, block)
		case token.CONTINUE:
			if entrance, exists := factory.loopEnds.Peek(); exists {
				factory.jumpTo(block, entrance)
			}
		case token.GOTO:
			label := cast.Label.Name

			// Backward declared label.
			if destination, exists := factory.labels[label]; exists {
				factory.connect(block, destination, nil)
			} else {
				// Forward declared label (Connection are deferred).
				if _, exists := factory.gotos[label]; !exists {
					factory.gotos[label] = []int{block}
				} else {
					factory.gotos[label] = append(factory.gotos[label], block)
				}

			}
		case token.FALLTHROUGH:
			factory.fallthroughs = append(factory.fallthroughs, block)
		}
	case *ast.LabeledStmt:
		front := factory.getBlock()
		label := cast.Label.Name
		block := factory.proceed(scope)
		if terminated := factory.isTerminated(front); !terminated {
			factory.connect(front, block, nil)
		}

		factory.statement(cast.Stmt)

		factory.labels[label] = block
		if gotos, exists := factory.gotos[label]; exists {
			for _, source := range gotos {
				factory.connect(source, block, nil)
			}
		}
	default:
		if terminated {
			return
		}

		factory.append(statement)
	}
}

func SCFG(function *ast.FuncDecl) *scfg.Graph[ast.Stmt, ast.Expr] {
	factory := new()
	factory.statement(function.Body)
	return factory.scopes
}
