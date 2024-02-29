package z3

import "github.com/Brandhoej/gobion/internal/z3"

type Config struct {
	_config *z3.Config
}

func NewConfig() *Config {
	return &Config{
		_config: z3.NewConfig(),
	}
}
