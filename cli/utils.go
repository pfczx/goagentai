package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func Exit() error {
	fmt.Println("Good bye!")
	os.Exit(0)
	return nil
}

func Help() error {
	commands := GetCommands()
	for _, cmd := range commands {
		fmt.Printf("Command: %s   Desc: %s\n", cmd.Name, cmd.Desc)
	}
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

func PrintState(profileName string, provider string, iternalProvider string, model string) {
	agent := color.New(color.FgCyan).Sprint("GoAgent")
	profile := color.New(color.FgCyan).Sprint(profileName)
	provider = color.New(color.FgGreen).Sprint(provider)
	iternalProvider = color.New(color.FgYellow).Sprint(iternalProvider)
	model = color.New(color.FgBlue).Sprint(model)
	prompt := fmt.Sprint("$")

	fmt.Printf("%s@%s (%s::%s::%s) \n", agent, profile, provider, iternalProvider, model)
	fmt.Printf("%s ", prompt)
}
