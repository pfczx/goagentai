package agent

import (
	"github.com/pfczx/goagentai/llm"
)

type Agent struct{
Provider llm.ModelProvider

}

func New(p llm.ModelProvider) *Agent {
	return &Agent{Provider: p}
}

func (a *Agent) Ask(input string) (string, error) {
	return a.Provider.Generate(input)
}
