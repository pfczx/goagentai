package cli

import (
	"fmt"

	"github.com/pfczx/goagentai/agent"
)

func SingleRun(agent *agent.Agent, args []string) {
	commandName := args[0]
	var cmdArgs []string
	if len(cmdArgs) > 1 {
		cmdArgs = args[1:]
	}
	err := HandleCommand(agent, commandName, cmdArgs...)
	if err != nil {
		fmt.Println("Error occured", err)
	}
}
