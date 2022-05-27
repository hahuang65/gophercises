package link_parser

import (
	"os"
	"reflect"
	"testing"
)

func TestEx1(t *testing.T) {
	doc, _ := os.Open("ex1.html")
	got := ParseLinks(doc)
	want := []Link{
		{
			Href: "/other-page",
			Text: "A link to another page",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx2(t *testing.T) {
	doc, _ := os.Open("ex2.html")
	got := ParseLinks(doc)
	want := []Link{
		{
			Href: "https://www.twitter.com/joncalhoun",
			Text: "Check me out on twitter",
		},
		{
			Href: "https://github.com/gophercises",
			Text: "Gophercises is on Github!",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx3(t *testing.T) {
	doc, _ := os.Open("ex3.html")
	got := ParseLinks(doc)
	want := []Link{
		{
			Href: "#",
			Text: "Login",
		},
		{
			Href: "/lost",
			Text: "Lost? Need help?",
		},
		{
			Href: "https://twitter.com/marcusolsson",
			Text: "@marcusolsson",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx4(t *testing.T) {
	doc, _ := os.Open("ex4.html")
	got := ParseLinks(doc)
	want := []Link{
		{
			Href: "/dog-cat",
			Text: "dog cat",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}
