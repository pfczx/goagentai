package cli

import (
	"fmt"
	"github.com/pfczx/goagentai/agent"
)

type CliCommand struct {
	Name     string
	Alias    string
	Desc     string
	Callback func(agent *agent.Agent, args ...string) error
}

func HandleCommand(agent *agent.Agent, commandName string, args ...string) error {
	commands := GetCommands()

	for _, cmd := range commands {
		if cmd.Alias == commandName {
			if err := cmd.Callback(agent, args...); err != nil {
				return err
			}
			return nil
		}
	}

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
			Name:  "exit",
			Alias: "e",
			Desc:  "exit repl",
			Callback: func(a *agent.Agent, _ ...string) error {
				return Exit(a)
			},
		},
		"help": {
			Name:  "help",
			Alias: "h",
			Desc:  "print description of commands",
			Callback: func(_ *agent.Agent, _ ...string) error {
				return Help()
			},
		},
		"init": {
			Name:  "init",
			Alias: "i",
			Desc:  "Creating profile with name provided in first argument and default config in .config/goagent",
			Callback: func(_ *agent.Agent, args ...string) error {
				return agent.InitProfile(args...)
			},
		},
		"ask": {
			Name:     "ask",
			Alias:    "a",
			Desc:     "Ask llm question provided in first argument",
			Callback: agent.RunAsk,
		},
		"switch": {
			Name:     "switch",
			Alias:    "s",
			Desc:     "Switching setting provided in first argument \n\nswitch [thing to swith] [name of new thing or alias] \n\nfirst arguments ->  profile : p,internal-provider : ip,provider : pv,model : m",
			Callback: agent.Switch,
		},
		"list": {
			Name:     "list",
			Alias:    "l",
			Desc:     "Print list of selected things \n\nlist [thing to list or alias] \n\nfirst arguments -> \n\nprofiles : p \n\nproviders : pv \n\nintenral-providers : ip \n\nmodels  : m --image flag for listing models accepting photo input",
			Callback: agent.List,
		},
		"config": {
			Name:     "config",
			Alias:    "c",
			Desc:     "Open config for current profile in default editor and load it after",
			Callback: agent.EditConfig,
		},
		"workspace": {
			Name:     "workspace",
			Alias:    "w",
			Desc:     "Execute workspace command provided in first argumet \n\nfirst arguments -> \n\ndd : a [filepath or directory] -p (p flag for persistent adding, without flag all content of workspace is erased after next ask)\n\nremove : r [filepath or directory] or all (clearing all workspace content) \n\nlist : l (listing all workspace content) all flag for listing with content",
			Callback: agent.Workspace,
		},
		"remove": {
			Name:     "remove",
			Alias:    "r",
			Desc:     "Removes specified thing in firts argument \n\nfirst arguments -> \n\nprofile : p [profile_name]\n\nhistory : h [history_type] (long-term : lt or short-term : st) \n\n used-tokens : ut (set used tokens to 0)",
			Callback: agent.Remove,
		},
	}
	return commands
}
