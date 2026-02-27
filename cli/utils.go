package cli

import (
	"fmt"
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
