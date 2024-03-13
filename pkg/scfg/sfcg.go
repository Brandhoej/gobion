package scfg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"strconv"
	"strings"
)

var (
	Initial  = -1
	Terminal = -2
)

type Block struct {
	statements []ast.Stmt
}

func NewBlock(statements ...ast.Stmt) *Block {
	if len(statements) == 1 && statements[0] == nil {
		return &Block{
			statements: nil,
		}
	}

	return &Block{
		statements: statements,
	}
}

func (block *Block) IsEmpty() bool {
	return len(block.statements) == 0
}

type Jump struct {
	expression ast.Expr
}

func NewJump(expression ast.Expr) *Jump {
	return &Jump{
		expression: expression,
	}
}

type Condition struct {
	expression ast.Expr
}

func NewCondition(expression ast.Expr) *Condition {
	return &Condition{
		expression: expression,
	}
}

type Graph struct {
	fset       *token.FileSet
	blocks     []*Block
	conditions []*Condition
	jumps      []*Jump
	// source -> condition -> destination -> jumps
	mapping map[int]map[int]map[int][]int
}

func New() *Graph {
	return &Graph{
		fset:       token.NewFileSet(),
		blocks:     make([]*Block, 0),
		conditions: make([]*Condition, 0),
		mapping:    map[int]map[int]map[int][]int{},
	}
}

func (scfg *Graph) DOT(writer io.Writer) {
	io.WriteString(writer, "digraph G {\n")

	blockName := func(index int) string {
		if index == Initial {
			return "INITIAL"
		} else if index == Terminal {
			return "TERMINAL"
		}
		return fmt.Sprintf("block_%v", index)
	}

	conditionName := func(condition int) string {
		return fmt.Sprintf("condition_%v", condition)
	}

	hasCondition := func(condition int) bool {
		return scfg.conditions[condition].expression != nil
	}

	astNodeStr := func(node any) string {
		if node == nil {
			return ""
		}
		var buffer bytes.Buffer
		format.Node(&buffer, scfg.fset, node)
		str := buffer.String()
		str = strings.ReplaceAll(str, "\"", "\"")
		str = strings.ReplaceAll(str, "\t", "")
		return str
	}

	attrs := func(mapping map[string]string) string {
		var builder strings.Builder
		counter := 0
		for key, value := range mapping {
			if counter > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s=%s", key, value))
			counter++
		}
		return builder.String()
	}

	dotNodeStr := func(name string, attributes map[string]string) string {
		return fmt.Sprintf("%s [%s]\n", name, attrs(attributes))
	}

	dotEdgeStr := func(source, target string, attributes map[string]string) string {
		return fmt.Sprintf("%s -> %s [%s]\n", source, target, attrs(attributes))
	}

	for idx, block := range scfg.blocks {
		name := blockName(idx)
		label := strconv.Quote(
			astNodeStr(block.statements),
		)
		io.WriteString(
			writer, dotNodeStr(name, map[string]string{
				"label": label,
				"shape": "rectangle",
			}),
		)
	}

	io.WriteString(
		writer,
		dotNodeStr(blockName(Initial), map[string]string{
			"shape": "point",
		}),
	)
	io.WriteString(
		writer,
		dotNodeStr(blockName(Terminal), map[string]string{
			"shape": "point",
		}),
	)

	for idx, condition := range scfg.conditions {
		if condition.expression != nil {
			name := conditionName(idx)
			label := strconv.Quote(
				astNodeStr(condition.expression),
			)
			io.WriteString(
				writer, dotNodeStr(name, map[string]string{
					"label": label,
					"shape": "diamond",
				}),
			)
		}
	}

	for source, conditions := range scfg.mapping {
		sourceName := blockName(source)

		for condition, destinations := range conditions {
			conditionName := conditionName(condition)

			if hasCondition(condition) {
				io.WriteString(
					writer, dotEdgeStr(sourceName, conditionName, nil),
				)
			}

			for destination, jumps := range destinations {
				destinationName := blockName(destination)

				for _, jump := range jumps {
					label := strconv.Quote(
						astNodeStr(scfg.jumps[jump].expression),
					)
					if hasCondition(condition) {
						io.WriteString(
							writer, dotEdgeStr(conditionName, destinationName, map[string]string{
								"label": label,
							}),
						)
					} else {
						io.WriteString(
							writer, dotEdgeStr(sourceName, destinationName, map[string]string{
								"label": label,
							}),
						)
					}
				}
			}
		}
	}

	// Close first scope subgraph cluster.
	io.WriteString(writer, "}\n")

	// Close digraph
	io.WriteString(writer, "}\n")
}
