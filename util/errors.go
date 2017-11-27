package util

import (
	"fmt"
	"os"
)

func Fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
