package agent

import (
	"fmt"
)

func Workspace(a *Agent, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("No subcommand provided")
	}

	switch args[0] {
	case "add", "a":
		if len(args) < 2 {
			return fmt.Errorf("Usage: workspace add [path] [-p]")
		}
		if len(args) > 3 {
			return fmt.Errorf("Too many arguments")
		}

		path := args[1]
		addOnce := true
		if len(args) == 3 && args[2] == "-p" {
			addOnce = false
		}

		err := a.WorspaceMenager.Add(path, addOnce)
		if err != nil {
			return err
		}

		fmt.Println("Added:", path)

	case "remove", "r":
		if len(args) < 2 {
			return fmt.Errorf("Usage: workspace remove [path or all]")
		}
		a.WorspaceMenager.Clear(args[1])

	case "list", "l":
		fmt.Println("Current content:")

		if len(args) >= 2 && args[1] == "all" {
			content := a.WorspaceMenager.WorkspaceToString(false)
			fmt.Println(content)
			return nil
		}

		for _, entry := range a.WorspaceMenager.Entries {
			fmt.Println(entry.Path)
		}
	default:
		return fmt.Errorf("First argument is not valid")
	}

	return nil
}
