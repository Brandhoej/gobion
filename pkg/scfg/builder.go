package scfg

import (
	"go/ast"
	"go/token"
)

var Unconditional = 0

type Builder struct {
	fset       *token.FileSet
	blocks     []*Block
	conditions []*Condition
	jumps      []*Jump
	// source -> condition -> destination -> jumps
	mapping map[int]map[int]map[int][]int
}

func NewBuilder(fset *token.FileSet) *Builder {
	builder := &Builder{
		fset:       fset,
		blocks:     []*Block{},
		conditions: []*Condition{{}}, // Initialize with a default empty condition
		jumps:      []*Jump{},
		mapping:    map[int]map[int]map[int][]int{},
	}

	// Add the initial and terminal nodes.
	builder.mapping[Initial] = map[int]map[int][]int{}
	builder.mapping[Terminal] = map[int]map[int][]int{}

	return builder
}

func (builder *Builder) AddBlock(block *Block) int {
	index := len(builder.blocks)
	builder.blocks = append(builder.blocks, block)
	builder.mapping[index] = map[int]map[int][]int{}
	return index
}

func (builder *Builder) Replace(index int, block *Block) {
	builder.blocks[index] = block
}

func (builder *Builder) Branch(source int, condition *Condition) int {
	index := len(builder.conditions)
	builder.conditions = append(builder.conditions, condition)
	builder.mapping[source][index] = map[int][]int{}
	return index
}

func (builder *Builder) Else(source, consequence int, alternative *Condition) int {
	// source -> consequence -> empty -> alternative.
	block := NewBlock(nil)
	empty := builder.AddBlock(block)
	builder.Contradiction(source, consequence, empty)
	return builder.Branch(empty, alternative)
}

func (builder *Builder) Condition(source, condition int) {
	builder.mapping[source][condition] = map[int][]int{}
}

func (builder *Builder) UnconditionalJump(source, destination int) int {
	return builder.ConditionalJump(
		source, Unconditional, destination, NewJump(nil),
	)
}

func (builder *Builder) Tautology(source, condition, destination int) int {
	return builder.ConditionalJump(source, condition, destination, NewJump(ast.NewIdent("true")))
}

func (builder *Builder) Contradiction(source, condition, destination int) int {
	return builder.ConditionalJump(source, condition, destination, NewJump(ast.NewIdent("false")))
}

func (builder *Builder) ConditionalJump(source, condition, destination int, jump *Jump) int {
	if builder.mapping[source][condition] == nil {
		builder.mapping[source][condition] = map[int][]int{}
	}

	jumps := builder.mapping[source][condition][destination]
	if jumps == nil {
		jumps = make([]int, 0, 1)
	}

	index := len(builder.jumps)
	builder.jumps = append(builder.jumps, jump)
	builder.mapping[source][condition][destination] = append(jumps, index)
	return index
}

func (builder *Builder) Build() *Graph {
	return &Graph{
		fset:       builder.fset,
		blocks:     builder.blocks,
		conditions: builder.conditions,
		jumps:      builder.jumps,
		mapping:    builder.mapping,
	}
}
