package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/memory"
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
	return NewAgent(profile, &memory.MemoryMenager{}), nil

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
	if resp.Usage != nil {
		fmt.Printf("Tokens prompt: %v completion: %v total: %v \n",
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
			resp.Usage.TotalTokens)
	} else {
		fmt.Println("(No token usage data available)")
	}

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

func List(agent *Agent, args ...string) error {
	var builder strings.Builder
	switch args[0] {
	case "providers":
		list := llm.ListProviders()
		builder.WriteString("# Currently availble providers\n\n")
		for _, provider := range list {
			builder.WriteString("## ")
			builder.WriteString(provider + "\n")
		}
		out, err := glamour.Render(builder.String(), "auto")
		if err != nil {
			return err
		}

		fmt.Print(out)

	case "iternal-providers":
		list, err := agent.Profile.Provider.ListIternalProviders()
		if err != nil {
			return err
		}
		builder.WriteString("# Currently availble iternal-providers for ")
		builder.WriteString(agent.Profile.Provider.Name() + "\n\n")
		for _, iternalProvider := range list {
			builder.WriteString("## ")
			builder.WriteString(iternalProvider + "\n")
		}
		out, err := glamour.Render(builder.String(), "auto")
		if err != nil {
			return err
		}

		fmt.Print(out)
	case "models":
		withPhoto := false
		if len(args) > 1 && args[1] == "--image" {
			withPhoto = true
		}
		list, err := agent.Profile.Provider.ListProviderModels(agent.Profile.Provider.IternalProviderName(), withPhoto)
		if err != nil {
			return err
		}
		builder.WriteString("# Currently availble models for ")
		builder.WriteString(agent.Profile.Provider.IternalProviderName() + "\n\n")
		for _, model := range list {
			builder.WriteString("## ")
			builder.WriteString(model + "\n")
		}
		out, err := glamour.Render(builder.String(), "auto")
		if err != nil {
			return err
		}

		fmt.Print(out)

	default:
		return fmt.Errorf("First argument is not valid")

	}
	return nil
}

func EditConfig(agent *Agent, args ...string) error {
	configPath := filepath.Join(agent.Profile.Path, "config.json")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	newAgent, err := InitAgent(agent.Profile.Name)
	if err != nil {
		return err
	}
	*agent = *newAgent
	return nil

}
