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

func Exit(msg string, statusCode int) {
	fmt.Println(msg)
	os.Exit(statusCode)
}

func Fail(msg string) {
	Exit(msg, 1)
}

func ReadCSV(path string) [][]string {
	csvFile, err := os.Open(path)
	defer csvFile.Close()

	if err != nil {
		Fail(fmt.Sprintf("Could not open CSV file `" + path + "`"))
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		Fail(fmt.Sprintf("Could not parse CSV file `" + path + "`"))
	}

	return lines
}

func ReadFile(path string) []byte {
	bytes, err := os.ReadFile(path)

	if err != nil {
		Fail(fmt.Sprintf("Could not read file `" + path + "`"))
	}

	return bytes
}
