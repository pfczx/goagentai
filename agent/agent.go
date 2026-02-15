package agent

import (
	"github.com/pfczx/goagentai/llm"
)

type Agent struct {
	Provider    llm.ModelProvider
	Temperature float64
	Config *Config
}

func NewAgent(p llm.ModelProvider, temp float64, c *Config) *Agent {
	return &Agent{
		Provider:    p,
		Temperature: temp,
		Config: c,
	}
}

func (a *Agent) Ask(input string) (string, error) {
	return a.Provider.Generate(input)
}
