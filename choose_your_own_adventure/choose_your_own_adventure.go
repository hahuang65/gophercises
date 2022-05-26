package choose_your_own_adventure

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"git.sr.ht/~hwrd/gophercises/exit"
)

type ChooseYourOwnAdventure struct{}

func (c *ChooseYourOwnAdventure) CommandName() string {
	return "choose_your_own_adventure"
}

func (c *ChooseYourOwnAdventure) Run(args []string) {
	var (
		jsonPath string
	)

	cmd := flag.NewFlagSet(c.CommandName(), flag.ExitOnError)
	cmd.StringVar(&jsonPath, "json", fmt.Sprintf("%s/gopher.json", c.CommandName()), "filepath to the YAML containing URLs")
	cmd.Parse(args)

	story := parseStory(jsonPath)
	html := template.Must(template.ParseFiles(fmt.Sprintf("%s/template.html", c.CommandName())))
	handler := storyHandler{story: story, template: html}

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	port := ":3000"

	fmt.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":3000", mux)
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

type storyHandler struct {
	story    story
	template *template.Template
}

func parseStory(path string) story {
	file, err := os.Open(path)

	if err != nil {
		exit.Fail("Could not open JSON file `" + path + "`")
	}

	defer file.Close()

	d := json.NewDecoder(file)
	var s story
	err = d.Decode(&s)

	if err != nil {
		exit.Fail("Could not unmarshal JSON")
	}

	return s
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.URL.Path)

	if p == "" || p == "/" {
		p = "/intro"
	}

	p = p[1:]

	if chapter, ok := h.story[p]; ok {
		err := h.template.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
