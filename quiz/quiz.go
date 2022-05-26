package quiz

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"git.sr.ht/~hwrd/gophercises/exit"
)

type Quiz struct{}

func (q *Quiz) CommandName() string {
	return "quiz"
}

func (q *Quiz) Run(args []string) {
	var (
		csvPath   string
		shuffle   bool
		timeLimit int
	)

	cmd := flag.NewFlagSet(q.CommandName(), flag.ExitOnError)
	cmd.StringVar(&csvPath, "csv", fmt.Sprintf("%s/problems.csv", q.CommandName()), "filepath to the CSV containing quiz questions")
	cmd.BoolVar(&shuffle, "shuffle", false, "Shuffle the order of the problems")
	cmd.IntVar(&timeLimit, "timer", 30, "Amount of time allowed")
	cmd.Parse(args)

	csvFile, err := os.Open(csvPath)
	defer csvFile.Close()

	if err != nil {
		exit.Fail(fmt.Sprintf("Could not open CSV file `" + csvPath + "`"))
	}

	problems := parseProblems(csvFile)
	poseProblems(problems, timeLimit, shuffle)
}

type problem struct {
	question string
	answer   string
}

func parseProblems(r io.Reader) []problem {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		exit.Fail("Could not parse CSV file.")
	}

	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func shuffleProblems(problems []problem) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func poseProblems(problems []problem, timeLimit int, shuffle bool) {
	correct := 0
	possible := len(problems)

	if shuffle {
		shuffleProblems(problems)
	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	go func() {
		<-timer.C
		fmt.Println() // Prints a newline before the exit message
		endQuiz(correct, possible)
	}()

	for _, problem := range problems {
		fmt.Printf("%s = ", problem.question)
		var answer string
		fmt.Scanf("%s\n", &answer)

		if checkAnswer(strings.TrimSpace(answer), problem) {
			correct++
		}
	}

	endQuiz(correct, possible)
}

func checkAnswer(a string, p problem) bool {
	return a == p.answer
}

func endQuiz(correct int, possible int) {
	exit.WithMsgAndCode(fmt.Sprintf("You answered %d correctly out of %d", correct, possible), 0)
}
