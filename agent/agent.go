package agent

import (
	"github.com/pfczx/goagentai/llm"
)

type Agent struct {
	provider    llm.ModelProvider
	temperature float64
}

func NewAgent(p llm.ModelProvider, temp float64) *Agent {
	return &Agent{
		provider:    p,
		temperature: temp,
	}
}


func (a *Agent) Ask(input string) (string, error) {
	return a.provider.Generate(input)
}
