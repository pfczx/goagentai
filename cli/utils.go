package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/pfczx/goagentai/agent"
)

func Exit(a *agent.Agent) error {
	fmt.Println("Good bye!")
	err := a.CloseDB()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
	return nil
}

func Help() error {
	commands := GetCommands()
	var builder strings.Builder

	builder.WriteString("# GoAgent Help\n\n")

	builder.WriteString("## Available Commands\n\n")

	for _, cmd := range commands {
		builder.WriteString(fmt.Sprintf("### `%s`\n\n", cmd.Name))
		builder.WriteString("Alias: " + cmd.Alias)
		builder.WriteString("\n\n")
		builder.WriteString(cmd.Desc)
		builder.WriteString("\n\n")
	}

	out, err := glamour.Render(builder.String(), "auto")
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}

func LoadEnv() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	envPath := filepath.Join(home, ".config", "goagent", ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("error reading .env from %s: %w", envPath, err)
	}

	return nil
}

func PrintState(profileName string, provider string, iternalProvider string, model string, tokenStatus string) {
	agent := color.New(color.FgCyan).Sprint("GoAgent")
	profile := color.New(color.FgCyan).Sprint(profileName)
	provider = color.New(color.FgGreen).Sprint(provider)
	iternalProvider = color.New(color.FgYellow).Sprint(iternalProvider)
	model = color.New(color.FgBlue).Sprint(model)
	prompt := fmt.Sprint("$")

	fmt.Printf("%s@%s (%s::%s::%s) %s\n", agent, profile, provider, iternalProvider, model, tokenStatus)
	fmt.Printf("%s ", prompt)
}
