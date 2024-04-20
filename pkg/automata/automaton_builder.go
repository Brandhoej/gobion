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
	guard Guard
	update Update
}

type EdgeConfiguration func(config *EdgeConfig)

func NewEdgeConfig(configs ...EdgeConfiguration) EdgeConfig {
	config := EdgeConfig{
		guard: NewTrueGuard(),
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
	initial graph.Key
	locations graph.Vertices[Location]
	edges graph.Edges[Edge]
}

func NewAutomatonBuilder() *AutomatonBuilder {
	return &AutomatonBuilder{
		locations: graph.NewVertexMap[Location](),
		edges: graph.NewEdgesMap[Edge](),
	}
}

func (builder *AutomatonBuilder) AddInitial(name string, configs ...LocationConfiguration) graph.Key {
	key := builder.AddLocation(name, configs...)
	builder.initial = key
	return key
}

func (builder *AutomatonBuilder) AddLocation(name string, configs ...LocationConfiguration) graph.Key {
	config := NewLocationConfig(configs...)
	location := NewLocation(name, config.invariant)
	key := builder.locations.Add(location)
	return key
}

func (builder *AutomatonBuilder) AddEdge(source, destination graph.Key, configs ...EdgeConfiguration) {
	config := NewEdgeConfig(configs...)
	edge := NewEdge(source, config.guard, config.update, destination)
	builder.edges.Add(edge)
}

func (builder *AutomatonBuilder) AddLoop(location graph.Key, configs ...EdgeConfiguration) {
	builder.AddEdge(location, location, configs...)
}

func (builder *AutomatonBuilder) Build(symbols symbols.Store[any]) Automaton {
	dg := graph.NewLabeledDirected[Edge, Location](
		builder.locations, builder.edges,
	)
	return *NewAutomaton(dg, builder.initial, symbols)
}