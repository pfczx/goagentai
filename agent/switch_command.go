package agent

import (
	"fmt"

	"github.com/pfczx/goagentai/llm"
)

func Switch(agent *Agent, args ...string) error {
	switch args[0] {
	case "profile", "p":
		newAgent, err := InitAgent(args[1])
		if err != nil {
			return err
		}
		*agent = *newAgent
		return nil

	case "provider", "pv":
		provider, err := llm.NewProvider(args[1],
			agent.Profile.Provider.ModelName(),
			agent.Profile.Provider.IternalProviderName())
		if err != nil {
			return err
		}
		agent.Profile.Provider = provider
	case "internal-provider", "ip":
		agent.Profile.Provider.SwitchIternalProvider(args[1])
	case "model", "m":
		agent.Profile.Provider.SwitchModel(args[1])
	default:
		return fmt.Errorf("First argument is not valid")
	}
	err := agent.Profile.UpdateConfigFromProfile(agent.MemoryMenager)
	if err != nil {
		return err
	}

	return nil
}
