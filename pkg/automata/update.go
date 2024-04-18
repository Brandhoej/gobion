package automata

import "github.com/Brandhoej/gobion/pkg/automata/language"

type Update struct {
	assignments []language.Assignment
}

func NewUpdate(assignments ...language.Assignment) Update {
	return Update{
		assignments: assignments,
	}
}

func (update Update) Apply(
	interpreter language.SymbolicStatementInterpreter, valuations language.Valuations,
) language.Valuations {
	copy := valuations.Copy()
	for idx := range update.assignments {
		interpreter.Assignment(update.assignments[idx])
	}
	return copy
}

func (update Update) String() string {
	printer := language.NewPrettyPrinter()
	for idx := range update.assignments {
		if idx > 0 {
			printer.Write(" ")
		}
		printer.Statement(update.assignments[idx])
	}
	return printer.String()
}