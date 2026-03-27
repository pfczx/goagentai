package agent

import (
	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/memory"
	"github.com/pfczx/goagentai/prompt"
	"github.com/pfczx/goagentai/workspace"
)

type Agent struct {
	Profile         *Profile
	MemoryMenager   *memory.MemoryMenager
	WorspaceMenager *workspace.WorkspaceMenager
}

func NewAgent(profile *Profile, memoryMenager *memory.MemoryMenager, workspaceMenager *workspace.WorkspaceMenager) *Agent {
	return &Agent{
		Profile:         profile,
		MemoryMenager:   memoryMenager,
		WorspaceMenager: workspaceMenager,
	}
}

func (a *Agent) Ask(input string, context []string) (*llm.ChatResponse, error) {
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

func (a *Agent) CloseDB() error {
	err := a.MemoryMenager.CloseDb()
	if err != nil {
		return err
	}
	return nil

}
