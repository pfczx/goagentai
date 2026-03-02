package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pfczx/goagentai/llm"
)

func InitAgent(profileName string) (*Agent, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ".config", "goagent", "profiles", profileName, "config.json")
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	profile, err := config.ProfileFromConfig()
	if err != nil {
		return nil, err
	}

	return NewAgent(profile), nil
}

func RunAsk(agent *Agent, args ...string) error {
	resp, err := agent.Ask(strings.Join(args, " "))
	if err != nil {
		return err
	}
	out, err := glamour.Render(resp.Text, "auto")
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil

}

func Switch(agent *Agent, args ...string) error {
	switch args[0] {
	case "profile":
		path, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(path, ".config", "goagent", "profiles", args[1], "config.json")
		conf, err := LoadConfig(path)
		if err != nil {
			return err
		}
		profile, err := conf.ProfileFromConfig()
		if err != nil {
			return err
		}
		agent.Profile = profile
		err = agent.Profile.SaveLatestUsedProfileName()
		if err != nil {
			return err
		}
	case "provider":
		provider, err := llm.NewProvider(args[1],
			agent.Profile.Provider.ModelName(),
			agent.Profile.Provider.IternalProviderName())
		if err != nil {
			return err
		}
		agent.Profile.Provider = provider
	case "iternal-provider":
		agent.Profile.Provider.SwitchIternalProvider(args[1])
	case "model":
		agent.Profile.Provider.SwitchModel(args[1])
	default:
		return fmt.Errorf("First argument is not valid")
	}
	err := agent.Profile.UpdateConfigFromProfile()
	if err != nil {
		return err
	}

	return nil
}
