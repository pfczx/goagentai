package main

import (
	"fmt"
	"github.com/pfczx/goagentai/agent"
	"github.com/pfczx/goagentai/cli"
	"os"
)

func main() {
	profileName, err := agent.LoadLatestUsedProfileName()
	if os.IsNotExist(err) {
		fmt.Println("No profile found. Running first initialization.")
		profileName = "default"
		agent.FirstInitialize()
	} else if err != nil {
		fmt.Println("Error loading latest profile: ", err)
		os.Exit(1)
	}
	err = cli.LoadEnv()
	if err != nil {
		fmt.Println("Error occured when loading env file: ", err)
	}
	agent, err := agent.InitAgent(profileName)
	if err != nil {
		fmt.Println("Error occured durning agent initialization: ", err)
		return
	}
	args := os.Args[1:]

	if len(args) == 0 {
		cli.Repl(agent)
		return
	}

	cli.SingleRun(agent, args)
}
