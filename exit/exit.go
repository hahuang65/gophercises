package exit

import (
	"fmt"
	"os"
)

func WithMsgAndCode(msg string, statusCode int) {
	fmt.Println(msg)
	os.Exit(statusCode)
}

func Fail(msg string) {
	WithMsgAndCode(msg, 1)
}
