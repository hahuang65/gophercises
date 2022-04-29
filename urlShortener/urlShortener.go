package urlShortener

import (
	"flag"
	"fmt"
	"net/http"

	"git.sr.ht/~hwrd/gophercises/util"
	"gopkg.in/yaml.v2"
)

type URLShortener struct{}

type url struct {
	Shortened string `yaml:"path"`
	Expanded  string `yaml:"url"`
}

func (u *URLShortener) CommandName() string {
	return "urlShortener"
}

func (u *URLShortener) Run(args []string) {
	// This exercise provided code.
	// `main.go` was copied into this function, after the `flag` code
	// `handler.go` had two empty functions, `MapHandler` and `YAMLHandler`,
	// both copied into this module, but outside the `Run` function.
	// I've tried to keep most of the sample code intact,
	// as the intention of the exercise is to simply implement the stubbed out
	// `MapHandler` and `YAMLHandler` functions.
	var (
		yamlPath string
	)

	cmd := flag.NewFlagSet(u.CommandName(), flag.ExitOnError)
	cmd.StringVar(&yamlPath, "yaml", fmt.Sprintf("%s/urls.yaml", u.CommandName()), "filepath to the YAML containing URLs")
	cmd.Parse(args)

	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// Build the MapHandler using the mux as the fallback
	mapHandler := mapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler := yamlHandler(yamlPath, mapHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func mapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requested := r.URL.Path

		if expanded, ok := pathsToUrls[requested]; ok {
			http.Redirect(w, r, expanded, http.StatusPermanentRedirect)
			fmt.Println("Found URL to expand!")
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func yamlHandler(yamlPath string, fallback http.Handler) http.HandlerFunc {
	mapped := mapFromYAML(yamlPath)
	return mapHandler(mapped, fallback)
}

func mapFromYAML(yamlPath string) map[string]string {
	var urls []url
	mappedYAML := make(map[string]string)

	yamlBytes := util.ReadFile(yamlPath)
	err := yaml.Unmarshal(yamlBytes, &urls)
	if err != nil {
		util.Fail("Could not parse YAML in `" + yamlPath + "`")
	}

	for _, url := range urls {
		mappedYAML[url.Shortened] = url.Expanded
	}

	return mappedYAML
}
