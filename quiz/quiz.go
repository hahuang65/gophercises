package quiz

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"git.sr.ht/~hwrd/gophercises/util"
)

type Quiz struct{}

func (q *Quiz) CommandName() string {
	return "quiz"
}

func (q *Quiz) Run(args []string) {
	var (
		csvPath string
		shuffle bool
	)

	cmd := flag.NewFlagSet(q.CommandName(), flag.ExitOnError)
	cmd.StringVar(&csvPath, "csv", fmt.Sprintf("%s/problems.csv", q.CommandName()), "filepath to the CSV containing quiz questions")
	cmd.BoolVar(&shuffle, "shuffle", false, "Shuffle the order of the problems")
	cmd.Parse(args)

	lines := util.ReadCSV(csvPath)
	problems := parseProblems(lines)
	score := poseProblems(problems, shuffle)

	fmt.Printf("Total questions: %d, Score: %d\n", len(problems), score)
}

type problem struct {
	question string
	answer   string
}

func parseProblems(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func poseProblems(problems []problem, shuffle bool) int {
	score := 0

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	for _, problem := range problems {
		fmt.Printf("%s = ", problem.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if strings.TrimSpace(answer) == problem.answer {
			score++
		}
	}

	return score
}
