package cli

import (
	"fmt"
	"os"
)

func Exit(args ...string) error {
	fmt.Println("Good bye!")
	os.Exit(0)
	return nil
}

func Help(args ...string) error {
	commands := GetCommands()
	for _, cmd := range commands {
		fmt.Printf("Command: %s   Desc: %s\n", cmd.Name, cmd.Desc)
	}
	return nil
}
