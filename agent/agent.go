package agent

import (
	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/memory"
	"github.com/pfczx/goagentai/prompt"
)

type Agent struct {
	Profile       *Profile
	MemoryMenager *memory.MemoryMenager
}

func NewAgent(profile *Profile, memoryMenager *memory.MemoryMenager) *Agent {
	return &Agent{
		Profile:       profile,
		MemoryMenager: memoryMenager,
	}
}

func (a *Agent) Ask(input string) (*llm.ChatResponse, error) {
	message, err := prompt.BuildAsk(input)
	if err != nil {
		return nil, err
	}
	llmResponse, err := a.Profile.Provider.Generate(message)
	if err != nil {
		return nil, err
	}
	return llmResponse, nil

}
