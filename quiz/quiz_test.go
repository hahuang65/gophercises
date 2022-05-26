package quiz

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestParseProblems(t *testing.T) {
	problems := `
5+5,10
2-1,1
3+4,7
`
	got := parseProblems(strings.NewReader(problems))
	want := []problem{
		{
			question: "5+5",
			answer:   "10",
		},
		{
			question: "2-1",
			answer:   "1",
		},
		{
			question: "3+4",
			answer:   "7",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
}

func TestShuffleProblems(t *testing.T) {
	problems := []problem{
		{
			question: "5+5",
			answer:   "10",
		},
		{
			question: "2-1",
			answer:   "1",
		},
		{
			question: "3+4",
			answer:   "7",
		},
	}

	shuffled := make([]problem, len(problems))
	copy(shuffled, problems)

	shuffleProblems(shuffled)

	if slices.Equal(problems, shuffled) {
		t.Errorf("\nproblems: %q\nis in the same order as\nshuffled: %q", problems, shuffled)
	}

	sort.Slice(problems, func(i, j int) bool {
		return problems[i].question < problems[j].question
	})
	sort.Slice(shuffled, func(i, j int) bool {
		return shuffled[i].question < shuffled[j].question
	})

	if !slices.Equal(problems, shuffled) {
		t.Errorf("\nproblems: %q\n does not contain the same contents as\nshuffled: %q", problems, shuffled)
	}
}

func TestCheckAnswer(t *testing.T) {
	p := problem{
		question: "2+3",
		answer:   "5",
	}

	got := checkAnswer("5", p)
	want := true

	if got != want {
		t.Error("\n checkAnswer returned `false`, but should be `true`")
	}

	got = checkAnswer("4", p)
	want = false

	if got != want {
		t.Error("\n checkAnswer returned `true`, but should be `false`")
	}
}
