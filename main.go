package main

import (
	"fmt"
	"os"

	"git.sr.ht/~hwrd/gophercises/quiz"
	"git.sr.ht/~hwrd/gophercises/util"
)

func addSubcommand(m map[string]util.Subcommand, cmd util.Subcommand) {
	m[cmd.CommandName()] = cmd
}

func main() {
	subcommands := make(map[string]util.Subcommand)
	addSubcommand(subcommands, &quiz.Quiz{})

	subcommand := os.Args[1]
	if cmd, ok := subcommands[subcommand]; ok {
		cmd.Run(os.Args[2:])
	} else {
		util.Fail(fmt.Sprintf("%s is not a valid subcommand", subcommand))
	}
}
