package cli

import "fmt"

type CliCommand struct {
	Name     string
	Desc     string
	Callback func(...string) error
}

func HandleCommand(commandName string, args ...string) error {
	commands := GetCommands()
	if cmd, exists := commands[commandName]; exists {
		if err := cmd.Callback(args...); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Command does not exists, type help for help :P")
	}
	return nil
}

func GetCommands() map[string]CliCommand {
	Commands := map[string]CliCommand{
		"exit": {
			Name:     "exit",
			Desc:     "exit pokedex",
			Callback: Exit,
		},
		"help": {
			Name:     "help",
			Desc:     "print desc of cmd",
			Callback: Help,
		},
	}
	return Commands
}
