package link_parser

import (
	"os"
	"reflect"
	"testing"
)

func TestEx1(t *testing.T) {
	doc, _ := os.Open("ex1.html")
	got := parseLinks(doc)
	want := []link{
		link{
			href: "/other-page",
			text: "A link to another page",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx2(t *testing.T) {
	doc, _ := os.Open("ex2.html")
	got := parseLinks(doc)
	want := []link{
		link{
			href: "https://www.twitter.com/joncalhoun",
			text: "Check me out on twitter",
		},
		link{
			href: "https://github.com/gophercises",
			text: "Gophercises is on Github!",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx3(t *testing.T) {
	doc, _ := os.Open("ex3.html")
	got := parseLinks(doc)
	want := []link{
		link{
			href: "#",
			text: "Login",
		},
		link{
			href: "/lost",
			text: "Lost? Need help?",
		},
		link{
			href: "https://twitter.com/marcusolsson",
			text: "@marcusolsson",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestEx4(t *testing.T) {
	doc, _ := os.Open("ex4.html")
	got := parseLinks(doc)
	want := []link{
		link{
			href: "/dog-cat",
			text: "dog cat",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}
