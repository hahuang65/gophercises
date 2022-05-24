package main

import (
	"fmt"
	"os"

	"git.sr.ht/~hwrd/gophercises/choose_your_own_adventure"
	"git.sr.ht/~hwrd/gophercises/link_parser"
	"git.sr.ht/~hwrd/gophercises/quiz"
	"git.sr.ht/~hwrd/gophercises/url_shortener"
	"git.sr.ht/~hwrd/gophercises/util"
)

func addSubcommand(m map[string]util.Subcommand, cmd util.Subcommand) {
	m[cmd.CommandName()] = cmd
}

func main() {
	subcommands := make(map[string]util.Subcommand)
	addSubcommand(subcommands, &choose_your_own_adventure.ChooseYourOwnAdventure{})
	addSubcommand(subcommands, &link_parser.LinkParser{})
	addSubcommand(subcommands, &quiz.Quiz{})
	addSubcommand(subcommands, &url_shortener.URLShortener{})

	subcommand := os.Args[1]
	if cmd, ok := subcommands[subcommand]; ok {
		cmd.Run(os.Args[2:])
	} else {
		util.Fail(fmt.Sprintf("%s is not a valid subcommand", subcommand))
	}
}
