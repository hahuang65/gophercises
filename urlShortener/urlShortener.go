package urlShortener

import (
	"flag"
	"fmt"
	"net/http"

	"git.sr.ht/~hwrd/gophercises/util"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

const BoltDBBucket = "URLShortener"

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
		dbPath   string
	)

	cmd := flag.NewFlagSet(u.CommandName(), flag.ExitOnError)
	cmd.StringVar(&yamlPath, "yaml", fmt.Sprintf("%s/urls.yaml", u.CommandName()), "filepath to the YAML containing URLs")
	cmd.StringVar(&dbPath, "db", fmt.Sprintf("%s/urls.db", u.CommandName()), "filepath to the BoltDB file containing URLs")
	cmd.Parse(args)

	mux := defaultMux()

	db, err := bolt.Open(dbPath, 0600, nil)
	defer db.Close()

	if err != nil {
		util.Fail("Could not open BoltDB file `" + dbPath + "`")
	}

	seedDB(db)
	if err != nil {
		util.Fail(fmt.Sprintf("Could not seed BoltDB: %s", err))
	}

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// Build the MapHandler using the mux as the fallback
	mapHandler := mapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler := yamlHandler(yamlPath, mapHandler)

	// Build the DBHandler using the YAMLhandler as the fallback
	dbHandler := dbHandler(db, yamlHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
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

func dbHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expanded []byte
		requested := r.URL.Path

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BoltDBBucket))
			expanded = b.Get([]byte(requested))

			return nil
		})

		if expanded != nil {
			http.Redirect(w, r, string(expanded), http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}

}

func seedDB(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BoltDBBucket))
		err = b.Put([]byte("/hwrd"), []byte("https://hwrd.me"))
		err = b.Put([]byte("/github"), []byte("https://github.com/hahuang65"))
		return err
	})

	return err
}
