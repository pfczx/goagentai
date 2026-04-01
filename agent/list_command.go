package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pfczx/goagentai/llm"
)

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
		for num, iternalProvider := range list {
			builder.WriteString("## ")
			builder.WriteString(fmt.Sprintf("%d ", num+1))
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
		for num, model := range list {
			builder.WriteString("## ")
			builder.WriteString(fmt.Sprintf("%d ", num+1))
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
