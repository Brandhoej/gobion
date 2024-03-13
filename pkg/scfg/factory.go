package scfg

import (
	"go/ast"
	"go/token"
)

type Sequence struct {
	blockIndex int
	sequenceIndex int
	block *Block
}

type Factory struct {
	sequences []*Sequence
	labels    map[string]*Sequence
	builder   *Builder
	pointer int
}

func NewFactory(fset *token.FileSet) *Factory {
	return &Factory{
		builder:   NewBuilder(fset),
		labels:    map[string]*Sequence{},
		sequences: []*Sequence{},
	}
}

func (factory *Factory) append(statement ast.Stmt) {
	sequence := factory.front()
	sequence.block.statements = append(sequence.block.statements, statement)
}

func (factory *Factory) front() *Sequence {
	return factory.sequences[factory.pointer]
}

func (factory *Factory) connect(source, destination int) {
	factory.builder.UnconditionalJump(source, destination)
}

func (factory *Factory) proceed() *Sequence {
	block := NewBlock()
	sequence := &Sequence{
		blockIndex: factory.builder.AddBlock(block),
		sequenceIndex: len(factory.sequences),
		block: block,
	}

	factory.sequences = append(factory.sequences, sequence)
	factory.pointer = sequence.sequenceIndex

	return sequence
}

func (factory *Factory) Branch(source int, expression ast.Expr) int {
	return factory.builder.Branch(source, NewCondition(expression))
}

func (factory *Factory) Merge(sequences ...*Sequence) *Sequence {
	var empty *Sequence = nil

	for _, current := range sequences {
		if current.block.IsEmpty() {
			empty = current
			break
		}
	}

	if empty != nil {
		for _, current := range sequences {
			if current == empty {
				continue
			}

			if current.block.IsEmpty() {
				factory.builder.Replace(current.blockIndex, empty.block)
			}
		}
	}

	return empty
}

func (factory *Factory) Combine(sequences ...*Sequence) (*Sequence, bool) {
	var proceeding *Sequence = factory.Merge(sequences...)
	
	merged := true
	if proceeding == nil {
		merged = false
		proceeding = factory.proceed()
	}

	for _, sequence := range sequences {
		if proceeding == sequence {
			continue
		}
		factory.builder.UnconditionalJump(
			sequence.blockIndex, proceeding.blockIndex,
		)
	}

	return proceeding, merged
}

func (factory *Factory) Consequence(source, condition int) *Sequence {
	consequence := factory.proceed()
	factory.builder.Tautology(
		source, condition, consequence.blockIndex,
	)
	return consequence
}

func (factory *Factory) Alternative(source, condition int) *Sequence {
	alternative := factory.proceed()
	factory.builder.Contradiction(
		source, condition, alternative.blockIndex,
	)
	return alternative
}

func (factory *Factory) Function(declaration *ast.FuncDecl) *Graph {
	sequence := factory.proceed()
	factory.Statement(declaration.Body)
	factory.builder.UnconditionalJump(
		Initial, sequence.blockIndex,
	)
	factory.builder.UnconditionalJump(
		factory.pointer, Terminal,
	)
	return factory.builder.Build()
}

func (factory *Factory) Statement(statement ast.Stmt) bool {
	switch cast := any(statement).(type) {
	case *ast.BlockStmt:
		for _, statement := range cast.List {
			if !factory.Statement(statement) {
				break
			}
		}
	case *ast.LabeledStmt:
		front := factory.front()
		proceeding := factory.proceed()
		factory.connect(front.blockIndex, proceeding.blockIndex)

		label := cast.Label.Name
		factory.labels[label] = proceeding

		factory.Statement(cast.Stmt)
	case *ast.BranchStmt:
		switch cast.Tok {
		case token.GOTO:
			front := factory.front()

			label := cast.Label.Name
			destination := factory.labels[label]

			source := factory.proceed()
			factory.connect(front.blockIndex, source.blockIndex)
			factory.append(cast)
			factory.connect(source.blockIndex, destination.blockIndex)

			following := factory.proceed()
			factory.connect(source.blockIndex, following.blockIndex)
		}
	case *ast.IfStmt:
		front := factory.front()
		condition := factory.Branch(front.blockIndex, cast.Cond)

		hasConsequence := len(cast.Body.List) > 0
		hasAlternative := cast.Else != nil && len(cast.Else.(*ast.BlockStmt).List) > 0

		if hasConsequence {
			factory.Consequence(front.blockIndex, condition)
			factory.Statement(cast.Body)
		}
		consequence := factory.front()

		if hasAlternative {
			factory.Alternative(front.blockIndex, condition)
			factory.Statement(cast.Else)
		}
		alternative := factory.front()

		if hasConsequence && hasAlternative {
			if junction, merged := factory.Combine(alternative, consequence); merged {
				factory.pointer = junction.sequenceIndex
			}
		} else {
			junction := factory.proceed()

			if hasConsequence {
				factory.builder.UnconditionalJump(consequence.blockIndex, junction.blockIndex)
			} else {
				factory.builder.Tautology(front.blockIndex, condition, junction.blockIndex)
			}
			
			if hasAlternative {
				factory.builder.UnconditionalJump(alternative.blockIndex, junction.blockIndex)
			} else {
				factory.builder.Contradiction(front.blockIndex, condition, junction.blockIndex)
			}
		}
	case *ast.ReturnStmt:
		front := factory.front()
		factory.append(statement)
		factory.builder.UnconditionalJump(
			front.blockIndex, Terminal,
		)
		factory.proceed()
		return false
	default:
		factory.append(statement)
	}
		
	return true
}
