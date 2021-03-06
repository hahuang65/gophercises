package link_parser

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"git.sr.ht/~hwrd/gophercises/exit"
	"golang.org/x/net/html"
)

type LinkParser struct{}

func (l *LinkParser) CommandName() string {
	return "link_parser"
}

func (l *LinkParser) Run(args []string) {
	var (
		htmlFile string
	)

	cmd := flag.NewFlagSet(l.CommandName(), flag.ExitOnError)
	cmd.StringVar(&htmlFile, "html", fmt.Sprintf("%s/ex1.html", l.CommandName()), "filepath to the HTML containing links")
	cmd.Parse(args)

	f, err := os.Open(htmlFile)
	if err != nil {
		exit.Fail("Could not open HTML file at " + htmlFile)
	}
	defer f.Close()

	fmt.Printf("Parsed links:\n%+v", ParseLinks(f))
}

type Link struct {
	Href string
	Text string
}

func ParseLinks(r io.Reader) []Link {
	doc, err := html.Parse(r)

	if err != nil {
		exit.Fail("Could not parse HTML")
	}

	nodes := linkNodes(doc)

	var links []Link
	for _, n := range nodes {
		links = append(links, linkify(n))
	}

	return links
}

func linkify(n *html.Node) Link {
	link := Link{}

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			link.Text = strings.Join(strings.Fields(linkText(n)), " ")
			break
		}
	}

	return link
}

func linkText(n *html.Node) string {
	if isText(n) {
		return n.Data
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += linkText(c)
	}
	return ret
}

func linkNodes(n *html.Node) []*html.Node {
	if isLink(n) {
		return []*html.Node{n}
	}

	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}

func isLink(node *html.Node) bool {
	if node.Type == html.ElementNode && node.Data == "a" {
		return true
	} else {
		return false
	}
}

func isText(node *html.Node) bool {
	if node.Type == html.TextNode {
		return true
	} else {
		return false
	}
}
