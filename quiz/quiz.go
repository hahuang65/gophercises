package quiz

import (
	"flag"
	"fmt"

	"git.sr.ht/~hwrd/gophercises/util"
)

type Quiz struct{}

func (q *Quiz) CommandName() string {
	return "quiz"
}

func (q *Quiz) Run(args []string) {
	var csvPath string

	cmd := flag.NewFlagSet(q.CommandName(), flag.ExitOnError)
	cmd.StringVar(&csvPath, "csv", fmt.Sprintf("%s/problems.csv", q.CommandName()), "filepath to the CSV containing quiz questions")
	cmd.Parse(args)

	lines := util.ReadCSV(csvPath)
	problems := parseProblems(lines)
	score := poseProblems(problems)

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
			answer:   line[1],
		}
	}

	return problems
}

func poseProblems(problems []problem) int {
	score := 0

	for _, problem := range problems {
		fmt.Printf("%s = ", problem.question)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == problem.answer {
			score++
		}
	}

	return score
}
