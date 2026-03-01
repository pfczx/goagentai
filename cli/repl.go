package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pfczx/goagentai/agent"
)

func Repl(agent *agent.Agent) error {
	sc := bufio.NewScanner(os.Stdin)
	for {
		PrintState(agent.Profile.Name,
			agent.Profile.Config.Provider,
			agent.Profile.Config.IternalProvider,
			agent.Profile.Config.Model,
		)
		if sc.Scan() {
			text := strings.TrimSpace(sc.Text())
			parts := strings.Fields(text)
			if len(parts) == 0 {
				continue
			}
			commandName := parts[0]
			var args []string
			if len(parts) > 1 {
				args = parts[1:]
			}
			err := HandleCommand(agent, commandName, args...)
			if err != nil {
				fmt.Println("Error occured: ", err)
			}
			if err := sc.Err(); err != nil {
				fmt.Println("Error reading input: ", err)
			}
		}
	}
}
