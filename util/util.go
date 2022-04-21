package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Subcommand interface {
	CommandName() string
	Run([]string)
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func ReadCSV(path string) [][]string {
	csvFile, err := os.Open(path)
	if err != nil {
		Exit(fmt.Sprintf("Could not read CSV file `" + path + "`"))
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		Exit(fmt.Sprintf("Could not parse CSV file `" + path + "`"))
	}

	return lines
}
