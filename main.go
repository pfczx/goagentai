package main

import (
	"os"

	"github.com/pfczx/goagentai/cli"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cli.Repl()
		return
	}

	cli.SingleRun(args)
}
