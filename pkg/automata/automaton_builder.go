package automata

import (
	"github.com/Brandhoej/gobion/pkg/graph"
	"github.com/Brandhoej/gobion/pkg/symbols"
)

type LocationConfig struct {
	invariant Invariant
}

type LocationConfiguration func(config *LocationConfig)

func NewLocationConfig(configs ...LocationConfiguration) LocationConfig {
	config := LocationConfig{
		invariant: NewTrueInvariant(),
	}
	for idx := range configs {
		configs[idx](&config)
	}
	return config
}

func WithInvariant(invariant Invariant) LocationConfiguration {
	return func(config *LocationConfig) {
		config.invariant = invariant
	}
}

type EdgeConfig struct {
	guard  Guard
	update Update
}

type EdgeConfiguration func(config *EdgeConfig)

func NewEdgeConfig(configs ...EdgeConfiguration) EdgeConfig {
	config := EdgeConfig{
		guard:  NewTrueGuard(),
		update: NewEmptyUpdate(),
	}
	for idx := range configs {
		configs[idx](&config)
	}
	return config
}

func WithGuard(guard Guard) EdgeConfiguration {
	return func(config *EdgeConfig) {
		config.guard = guard
	}
}

func WithUpdate(update Update) EdgeConfiguration {
	return func(config *EdgeConfig) {
		config.update = update
	}
}

type AutomatonBuilder struct {
	initial   symbols.Symbol
	locations graph.Vertices[symbols.Symbol, Location]
	edges     graph.Edges[symbols.Symbol, Edge]
	factory   *symbols.SymbolsFactory
}

func NewAutomatonBuilder() *AutomatonBuilder {
	return &AutomatonBuilder{
		locations: graph.NewVertexMap[symbols.Symbol, Location](),
		edges:     graph.NewEdgesMap[symbols.Symbol, Edge](),
		factory:   symbols.NewSymbolsFactory(),
	}
}

func (builder *AutomatonBuilder) AddInitial(name string, configs ...LocationConfiguration) symbols.Symbol {
	key := builder.AddLocation(name, configs...)
	builder.initial = key
	return key
}

func (builder *AutomatonBuilder) AddLocation(name string, configs ...LocationConfiguration) symbols.Symbol {
	config := NewLocationConfig(configs...)
	location := NewLocation(name, config.invariant)
	symbol := symbols.Symbol(builder.factory.Next())
	key := builder.locations.Add(location, symbol)
	return key
}

func (builder *AutomatonBuilder) AddEdge(source, destination symbols.Symbol, configs ...EdgeConfiguration) {
	config := NewEdgeConfig(configs...)
	edge := NewEdge(source, config.guard, config.update, destination)
	builder.edges.Connect(edge)
}

func (builder *AutomatonBuilder) AddLoop(location symbols.Symbol, configs ...EdgeConfiguration) {
	builder.AddEdge(location, location, configs...)
}

func (builder *AutomatonBuilder) Build() SymbolicAutomaton {
	dg := graph.NewLabeledDirected[symbols.Symbol, Edge, Location](
		builder.locations, builder.edges,
	)
	return *NewSymbolicAutomaton(*NewAutomaton(dg, builder.initial))
}
