package cli

import (
	"fmt"
	"github.com/pfczx/goagentai/agent"
)

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
	commands := map[string]CliCommand{
		"exit": {
			Name:     "exit",
			Desc:     "exit repl",
			Callback: Exit,
		},
		"help": {
			Name:     "help",
			Desc:     "print description of commands",
			Callback: Help,
		},
		"init": {
			Name:     "init",
			Desc:     "Creating profile with name provided in first argument and default config in .config/goagent",
			Callback: agent.InitProfile,
		},
	}
	return commands
}
