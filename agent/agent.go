package agent

import (
	"fmt"

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

func (a *Agent) Ask(input string) (*llm.ChatResponse, error) {
	var context []string
	if a.MemoryMenager.MemoryOn {
		//clearing workspace entries set to addOnce
		workspaceString := a.WorspaceMenager.WorkspaceToString(true)
		//last n conversations
		shortTermString := a.MemoryMenager.ShortTermToString()
		//prompt + short term used for searching relevant long term chunks
		longTermString, err := a.MemoryMenager.GetLongTermContextString(fmt.Sprintf("%s %s", input, shortTermString))

		if err != nil {
			return nil, err

		}
		if workspaceString != "" {
			context = append(context, workspaceString)
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
