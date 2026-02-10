package cli

import (
	"fmt"
)

func SingleRun(args []string) {
	commandName := args[0]
	var cmdArgs []string
	if len(cmdArgs) > 1 {
		cmdArgs = args[1:]
	}
	err := HandleCommand(commandName, cmdArgs...)
	if err != nil {
		fmt.Println("Error occured", err)
	}
}
