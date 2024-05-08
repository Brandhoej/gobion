package automata

import (
	"fmt"

	"github.com/Brandhoej/gobion/pkg/symbols"
)

type Action symbols.Symbol

func (action Action) String(input bool, store symbols.Store[any]) string {
	name, _ := store.Item(symbols.Symbol(action))
	if input {
		return fmt.Sprintf("%s?", name)
	}
	return fmt.Sprintf("%s!", name)
}