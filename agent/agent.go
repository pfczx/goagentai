package agent

import (
	"fmt"

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
		shortTermString := a.MemoryMenager.ShortTermToString()
		//prompt + short term used for searching relevant long term chunks
		longTermString, err := a.MemoryMenager.GetLongTermContextString(fmt.Sprintf("%s %s", input, shortTermString))
		if err != nil {
			return nil, err

		}
		if shortTermString != "" {
			context = append(context, shortTermString)
		}
		if longTermString != "" {
			context = append(context, longTermString)
		}

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

func (a *Agent) CloseDB() error {
	err := a.MemoryMenager.CloseDb()
	if err != nil {
		return err
	}
	return nil

}
