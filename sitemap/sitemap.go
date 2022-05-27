package sitemap

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"git.sr.ht/~hwrd/gophercises/exit"
	"git.sr.ht/~hwrd/gophercises/link_parser"
)

type Sitemap struct{}

func (s *Sitemap) CommandName() string {
	return "sitemap"
}

func (s *Sitemap) Run(args []string) {
	var (
		site string
	)

	cmd := flag.NewFlagSet(s.CommandName(), flag.ExitOnError)
	cmd.StringVar(&site, "site", "calhoun.io", "URL of the site to map")
	cmd.Parse(args)

	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, link := range mapSite(site) {
		toXml.Urls = append(toXml.Urls, loc{link})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		exit.Fail(fmt.Sprintf("Couldn't encode to XML: %s", err))
	}
	fmt.Println()
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func isLocalLink(l link_parser.Link, domain string) bool {
	href := l.Href

	if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "http") && domainMatches(href, domain) {
		return true
	} else {
		return false
	}

}

func domainMatches(href string, domain string) bool {
	parsedHref, err := url.Parse(ensureScheme(href))
	if err != nil {
		exit.Fail(fmt.Sprintf("Could not parse URL %q", href))
	}

	hostname := removeWWW(parsedHref.Host)
	domain = removeWWW(domain)

	if strings.HasPrefix(hostname, domain) {
		return true
	} else {
		return false
	}
}

func ensureScheme(s string) string {
	if !strings.HasPrefix(s, "http") {
		s = "http://" + s
	}

	return s
}

func removeWWW(s string) string {
	if strings.HasPrefix(s, "www.") {
		s = strings.Replace(s, "www.", "", 1)
	}

	return s
}

func mapPage(r io.Reader, domain string) []link_parser.Link {
	var ret []link_parser.Link

	for _, l := range link_parser.ParseLinks(r) {
		if isLocalLink(l, domain) {
			ret = append(ret, ensureDomain(l, domain))
		}
	}

	return ret
}

func ensureDomain(l link_parser.Link, domain string) link_parser.Link {
	href := removeWWW(l.Href)
	domain = removeWWW(domain)

	parsedHref, err := url.Parse(ensureScheme(href))
	if err != nil {
		exit.Fail(fmt.Sprintf("Could not parse URL %q", href))
	}

	if parsedHref.Hostname() == "" {
		l.Href = domain + parsedHref.Path
	}

	return l
}

func mapSite(u string) []string {
	seen := newStringSet()
	queue := newStringSet()
	nextQueue := newStringSet(u)

	// Arbitrary depth to map the site
	for i := 0; i < 3; i++ {
		queue, nextQueue = nextQueue, newStringSet()

		for x := range queue {
			if seen.include(x) {
				continue
			}

			seen.add(x)
			r, err := http.Get(ensureScheme(x))
			if err != nil {
				exit.Fail(fmt.Sprintf("Could not fetch %q", x))
			}
			defer r.Body.Close()

			for _, l := range mapPage(r.Body, u) {
				nextQueue.add(ensureScheme(removeWWW(l.Href)))
			}
		}
	}

	return seen.members()
}

type stringSet map[string]struct{}
type empty struct{}

func (s stringSet) add(member string) {
	s[member] = empty{}
}

func newStringSet(members ...string) stringSet {
	ret := stringSet{}

	for _, m := range members {
		ret.add(m)
	}

	return ret
}

func (s stringSet) len() int {
	return len(s)
}

func (s stringSet) include(key string) bool {
	if _, ok := s[key]; ok {
		return true
	} else {
		return false
	}
}

func (s stringSet) members() []string {
	ret := make([]string, 0, s.len())

	for k, _ := range s {
		ret = append(ret, k)
	}

	return ret

}
