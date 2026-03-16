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
	var context []string
	if a.MemoryMenager.MemoryOn {
		context = append(context, a.MemoryMenager.ShortTermToString())
	}
	message, err := prompt.BuildAsk(input, context)
	if err != nil {
		return nil, err
	}
	llmResponse, err := a.Profile.Provider.Generate(message)
	if err != nil {
		return nil, err
	}
	return llmResponse, nil

}
