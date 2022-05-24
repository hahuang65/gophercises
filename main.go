package main

import (
	"fmt"
	"os"

	"git.sr.ht/~hwrd/gophercises/chooseYourOwnAdventure"
	"git.sr.ht/~hwrd/gophercises/linkParser"
	"git.sr.ht/~hwrd/gophercises/quiz"
	"git.sr.ht/~hwrd/gophercises/urlShortener"
	"git.sr.ht/~hwrd/gophercises/util"
)

func addSubcommand(m map[string]util.Subcommand, cmd util.Subcommand) {
	m[cmd.CommandName()] = cmd
}

func main() {
	subcommands := make(map[string]util.Subcommand)
	addSubcommand(subcommands, &chooseYourOwnAdventure.ChooseYourOwnAdventure{})
	addSubcommand(subcommands, &linkParser.LinkParser{})
	addSubcommand(subcommands, &quiz.Quiz{})
	addSubcommand(subcommands, &urlShortener.URLShortener{})

	subcommand := os.Args[1]
	if cmd, ok := subcommands[subcommand]; ok {
		cmd.Run(os.Args[2:])
	} else {
		util.Fail(fmt.Sprintf("%s is not a valid subcommand", subcommand))
	}
}
