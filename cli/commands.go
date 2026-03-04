package cli

import (
	"fmt"
	"github.com/pfczx/goagentai/agent"
)

type CliCommand struct {
	Name     string
	Desc     string
	Callback func(agent *agent.Agent, args ...string) error
}

func HandleCommand(agent *agent.Agent, commandName string, args ...string) error {
	commands := GetCommands()
	if cmd, exists := commands[commandName]; exists {
		if err := cmd.Callback(agent, args...); err != nil {
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
			Name: "exit",
			Desc: "exit repl",
			Callback: func(_ *agent.Agent, _ ...string) error {
				return Exit()
			},
		},
		"help": {
			Name: "help",
			Desc: "print description of commands",
			Callback: func(_ *agent.Agent, _ ...string) error {
				return Help()
			},
		},
		"init": {
			Name: "init",
			Desc: "Creating profile with name provided in first argument and default config in .config/goagent",
			Callback: func(_ *agent.Agent, args ...string) error {
				return agent.InitProfile(args...)
			},
		},
		"ask": {
			Name:     "ask",
			Desc:     "Ask llm question provided in first argument",
			Callback: agent.RunAsk,
		},
		"switch": {
			Name:     "switch",
			Desc:     "Switching setting provided in first argument, switch [thing to swith] [name of new thing], first arguments: profile,iternal-provider,provider,model",
			Callback: agent.Switch,
		},
		"list": {
			Name:     "list",
			Desc:     "Print list of selected things, list [thing to list],first arguments: providers,itenral-providers,models ,--image flag for listing models accepting photo input",
			Callback: agent.List,
		},
		"config": {
			Name:     "config",
			Desc:     "Open config for current profile in default editor and load it after",
			Callback: agent.EditConfig,
		},
	}
	return commands
}
