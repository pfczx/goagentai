package agent

import (
	"fmt"
	"strconv"

	"github.com/pfczx/goagentai/llm"
)

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func Switch(agent *Agent, args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing sub‑command or argument")
	}
	switch args[0] {
	case "profile", "p":
		newAgent, err := InitAgent(args[1])
		if err != nil {
			return err
		}
		*agent = *newAgent
		agent.Profile.SaveLatestUsedProfileName()
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
		if isNumber(args[1]) {
			providers, err := agent.Profile.Provider.
				ListIternalProviders()
			if err != nil {
				return fmt.Errorf("error fetching internal providers: %w",
					err)
			}

			num, _ := strconv.Atoi(args[1])
			if num < 1 || num > len(providers) {
				return fmt.Errorf("internal provider number out of range  (1-%d)", len(providers))
			}
			agent.Profile.Provider.SwitchIternalProvider(providers[num-1])
		} else {
			agent.Profile.Provider.SwitchIternalProvider(args[1])
		}

	case "model", "m":
		if isNumber(args[1]) {
			var models []string
			var err error
			if len(args) == 3 && args[2] == "-img" {
				models, err = agent.Profile.Provider.ListProviderModels(
					agent.Profile.Provider.IternalProviderName(), true)
				if err != nil {
					return fmt.Errorf("error fetching models: %w", err)
				}

			} else {
				models, err = agent.Profile.Provider.ListProviderModels(
					agent.Profile.Provider.IternalProviderName(), false)
				if err != nil {
					return fmt.Errorf("error fetching models: %w", err)
				}

			}

			num, _ := strconv.Atoi(args[1])
			if num < 1 || num > len(models) {
				return fmt.Errorf("model number out of range (1-%d)",
					len(models))
			}
			agent.Profile.Provider.SwitchModel(models[num-1])
		} else {
			agent.Profile.Provider.SwitchModel(args[1])
		}

	default:
		return fmt.Errorf("First argument is not valid")
	}
	err := agent.Profile.UpdateConfigFromProfile(agent.MemoryMenager)
	if err != nil {
		return err
	}

	return nil
}
