package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

func Remove(a *Agent, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("No subcommand provided")
	}

	switch args[0] {
	case "profile", "p":
		if a.Profile.Name == args[1] {
			return fmt.Errorf("Switch profile before deleting")
		}
		if !ProfileExists(args[1]) {
			return fmt.Errorf("profile does not exist")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path := filepath.Join(homeDir, ".config", "goagent", "profiles", args[1])

		if err := os.RemoveAll(path); err != nil {
			return err
		}
		err = a.MemoryMenager.ClearLongTermMemory(args[1])
		if err != nil {
			return err
		}

		return nil

	case "history", "h":
		if args[1] == "short-term" || args[1] == "st" {
			a.MemoryMenager.ClearShortTerm()
			return nil

		} else if args[1] == "long-term" || args[1] == "lt" {
			err := a.MemoryMenager.ClearLongTermMemory(a.Profile.Name)
			if err != nil {
				return err
			}
			return nil
		} else {
			return fmt.Errorf("Wrong history type argument")
		}

	}
	return fmt.Errorf("Wrong first argument")
}
