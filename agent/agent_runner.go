package agent

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/memory"
	"github.com/pfczx/goagentai/workspace"
)

func InitMemoryMenagerFromConfig(config *Config) (*memory.MemoryMenager, error) {
	menager, err := memory.InitMenager(
		config.Path,
		config.Name,
		config.MemoryOn,
		config.ShortTermMemoryLimit,
		config.ShortTermMemoryEvaluation,
		config.LongTermMemoryBufferSize,
		config.LongTermMemoryChunkSize,
		config.LongTermMemoryStorageSize,
		config.LongTermMemoryChunksToAdd,
		config.LongTermMemorySummarizationProvider,
		config.LongTermMemorySummarizationInternalProvider,
		config.LongTermMemorySummarizationModel,
	)

	if err != nil {
		return nil, err
	}
	return menager, nil

}

func InitAgent(profileName string) (*Agent, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ".config", "goagent", "profiles", profileName, "config")
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	profile, err := config.ProfileFromConfig()
	if err != nil {
		return nil, err
	}
	Memorymenager, err := InitMemoryMenagerFromConfig(config)
	if err != nil {
		return nil, err
	}
	WorkspaceMenager := workspace.NewWorkspaceMenager()

	return NewAgent(profile, Memorymenager, WorkspaceMenager), nil

}

func RunAsk(agent *Agent, args ...string) error {
	prompt := strings.Join(args, " ")
	resp, err := agent.Ask(prompt)
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
		if agent.MemoryMenager.MemoryOn {
			usefull := false
			if agent.MemoryMenager.ShortTermMemoryEvaluation {
				sc := bufio.NewScanner(os.Stdin)
				fmt.Println("Was it usefull? [y/anything]")

				if sc.Scan() {
					if sc.Text() == "y" {
						usefull = true
					}
				}
			}
			agent.MemoryMenager.AppendShortTermHistory(prompt, resp.Text, usefull)
			if err = agent.MemoryMenager.UpdateLongTerm(); err != nil {
				return err
			}
		}

	} else {
		fmt.Println("no token usage error : RunAsk")
	}

	return nil
}

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

func List(agent *Agent, args ...string) error {
	var builder strings.Builder
	switch args[0] {
	case "profiles", "p":
		path, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(path, ".config", "goagent", "profiles")

		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("cannot read profiles directory: %w", err)
		}

		builder.WriteString("# Available profiles\n\n")
		for _, entry := range entries {
			if entry.IsDir() {
				builder.WriteString("## " + entry.Name() + "\n")
			}
		}

		out, err := glamour.Render(builder.String(), "auto")
		if err != nil {
			return err
		}
		fmt.Print(out)
	case "providers", "pv":
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

	case "internal-providers", "ip":
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
	case "models", "m":
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
	configPath := filepath.Join(agent.Profile.Path, "config")
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
