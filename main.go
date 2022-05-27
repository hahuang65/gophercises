package main

import (
	"fmt"
	"os"

	"git.sr.ht/~hwrd/gophercises/choose_your_own_adventure"
	"git.sr.ht/~hwrd/gophercises/exit"
	"git.sr.ht/~hwrd/gophercises/link_parser"
	"git.sr.ht/~hwrd/gophercises/quiz"
	"git.sr.ht/~hwrd/gophercises/sitemap"
	"git.sr.ht/~hwrd/gophercises/url_shortener"
)

type subcommand interface {
	CommandName() string
	Run([]string)
}

func addSubcommand(m map[string]subcommand, cmd subcommand) {
	m[cmd.CommandName()] = cmd
}

func main() {
	subcommands := make(map[string]subcommand)
	addSubcommand(subcommands, &choose_your_own_adventure.ChooseYourOwnAdventure{})
	addSubcommand(subcommands, &link_parser.LinkParser{})
	addSubcommand(subcommands, &quiz.Quiz{})
	addSubcommand(subcommands, &url_shortener.URLShortener{})
	addSubcommand(subcommands, &sitemap.Sitemap{})

	subcommand := os.Args[1]
	if cmd, ok := subcommands[subcommand]; ok {
		cmd.Run(os.Args[2:])
	} else {
		exit.Fail(fmt.Sprintf("%s is not a valid subcommand", subcommand))
	}
}
