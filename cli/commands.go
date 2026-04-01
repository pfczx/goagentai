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
			Desc:     "Ask [flag] [prompt] (spaces works in prompt) llm question provided in first argument, -s flag for adding screenshot with github.com/kbinani/screenshot, Your system should meet the requirements of dependencies if -s flag is used: \n\nLinux/FreeBSD - https://github.com/BurntSushi/xgb \n\nOSX - cgo. \n\nDo not panic if errors occur, simply check if model see your screen with simple prompt \n\nPhotos are converted to base64 format, it may be too big for low token limit api",
			Callback: agent.RunAsk,
		},
		"switch": {
			Name:     "switch",
			Alias:    "s",
			Desc:     "Switching setting provided in first argument, instead of a name you can provide numbers from command list without flags \n\nswitch [thing to swith or alias] [name of new thing ] \n\nfirst arguments ->  profile : p  internal-provider : ip  provider : pv  model : m",
			Callback: agent.Switch,
		},
		"list": {
			Name:     "list",
			Alias:    "l",
			Desc:     "Print list of selected things \n\nlist [thing to list or alias] \n\nfirst arguments -> profiles : p  providers : pv  intenral-providers : ip  models  : m -img flag for listing models accepting photo input(only works for HuggingFace, for other providers you need to visit their website )",
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
			Desc:     "Execute workspace command provided in first argumet \n\nfirst arguments ->  add : a [filepath or directory] -p (p flag for persistent adding, without flag all content of workspace is erased after next ask)  remove : r [filepath or directory] or all (clearing all workspace content)  list : l (listing all workspace paths) all flag for listing with content",
			Callback: agent.Workspace,
		},
		"remove": {
			Name:     "remove",
			Alias:    "r",
			Desc:     "Removes specified thing in firts argument \n\nfirst arguments ->  profile : p [profile_name]  history : h [history_type] (long-term : lt or short-term : st)  used-tokens : ut (set used tokens to 0)",
			Callback: agent.Remove,
		},
	}
	return commands
}
