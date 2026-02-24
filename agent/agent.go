package agent

import (
	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/prompt"
)

type Agent struct {
	profile *Profile
}

func NewAgent(profile *Profile) *Agent {
	return &Agent{
		profile: profile,
	}
}

func (a *Agent) Ask(input string) (*llm.ChatResponse, error) {
	message, err := prompt.BuildAsk(input)
	if err != nil {
		return nil, err
	}
	llmResponse, err := a.profile.Provider.Generate(message)
	if err != nil {
		return nil, err
	}
	return llmResponse, nil

}
