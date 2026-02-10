package cli

import(
"os"
	"bufio"
	"fmt"
	"strings"
)

func Repl() error {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("\033[32m" + "GoAgent > " + "\033[0m")
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
			err := HandleCommand(commandName,args...)
			if err!=nil{
				fmt.Println("Error occured",err)
			}
			if err := sc.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		}
	}
}
