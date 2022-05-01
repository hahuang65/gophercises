package chooseYourOwnAdventure

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"git.sr.ht/~hwrd/gophercises/util"
)

type ChooseYourOwnAdventure struct{}

func (c *ChooseYourOwnAdventure) CommandName() string {
	return "chooseYourOwnAdventure"
}

func (c *ChooseYourOwnAdventure) Run(args []string) {
	var (
		jsonPath string
	)

	cmd := flag.NewFlagSet(c.CommandName(), flag.ExitOnError)
	cmd.StringVar(&jsonPath, "json", fmt.Sprintf("%s/gopher.json", c.CommandName()), "filepath to the YAML containing URLs")
	cmd.Parse(args)

	story := parseStory(jsonPath)
	fmt.Println(story)
}

type story map[string]chapter

type chapter struct {
	Title   string   `json:"title"`
	Text    []string `json:"story"`
	Choices []choice `json:"options"`
}

type choice struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func parseStory(path string) story {
	file, err := os.Open(path)

	if err != nil {
		util.Fail("Could not open JSON file `" + path + "`")
	}

	defer file.Close()

	d := json.NewDecoder(file)
	var s story
	err = d.Decode(&s)

	if err != nil {
		util.Fail("Could not unmarshal JSON")
	}

	return s
}
