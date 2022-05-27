package sitemap

import (
	"os"
	"testing"

	"git.sr.ht/~hwrd/gophercises/link_parser"
	"golang.org/x/exp/slices"
)

func TestIsLocalLink(t *testing.T) {
	cases := []struct {
		name   string
		href   string
		domain string
		want   bool
	}{
		{"WithRelativePath", "/bar", "hwrd.me", true},
		{"WithScheme", "https://hwrd.me/foo", "hwrd.me", true},
		{"WithMismatchedDomain", "https://hwrd.me/foo", "foo.com", false},
		{"WithoutScheme", "facebook.com/", "facebook.com", true},
		{"WithWWW", "www.facebook.com", "facebook.com", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := link_parser.Link{
				Href: c.href,
			}

			got := isLocalLink(l, c.domain)
			if got != c.want {
				t.Errorf("got: %t, expected: %t", got, c.want)
			}
		})
	}
}

func TestMapPageOnlyReturnsLocalLinks(t *testing.T) {
	doc, _ := os.Open("test.html")

	got := mapPage(doc, "twitter.com")
	want := []link_parser.Link{
		{
			Href: "https://www.twitter.com/joncalhoun",
			Text: "Check me out on twitter",
		},
		{
			Href: "twitter.com/foobar",
			Text: "Foo Bar",
		},
	}

	if !slices.Equal(got, want) {
		t.Errorf("\n	got: %+v\nwant: %+v", got, want)
	}
}

func TestEnsureDomain(t *testing.T) {
	cases := []struct {
		name   string
		href   string
		domain string
		want   string
	}{
		{"AddsDomainWhenNoneExists", "/foo", "bar.com", "bar.com/foo"},
		{"DoesNothingWhenDomainExists", "bar.com/foo", "bar.com", "bar.com/foo"},
		{"DoesNothingWhenWWWExists", "www.bar.com/foo", "bar.com", "www.bar.com/foo"},
		{"DoesNothingDomainDoesNotMatch", "baz.com/foo", "bar.com", "baz.com/foo"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := link_parser.Link{
				Href: c.href,
			}

			got := ensureDomain(l, c.domain).Href

			if got != c.want {
				t.Errorf("got: %q, expected: %q", got, c.want)
			}
		})
	}
}
