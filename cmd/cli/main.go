package main

import (
	"github.com/andrewbackes/tourney/cmd/cli/instructions"
	"github.com/andrewbackes/tourney/util"
	"os"
	"strings"
)

func command() (verb, noun, subject string) {
	var cmd []string
	for i, v := range os.Args {
		if i != 0 && !strings.HasPrefix(v, "--") {
			cmd = append(cmd, v)
		}
	}
	if len(cmd) < 3 {
		util.Fail("Not enough arguements")
	}
	return cmd[0], cmd[1], cmd[2]
}

func flags() map[string]string {
	a := map[string]string{}
	for _, v := range os.Args {
		if strings.HasPrefix(v, "--") {
			key := strings.Split(v, "=")[0][2:]
			val := strings.Split(v, "=")[1]
			a[key] = val
		}
	}
	return a
}

func main() {
	verb, noun, subject := command()
	flags := flags()

	if verb == "add" && noun == "engine" {
		if _, exists := flags["filepath"]; !exists {
			flags["filepath"] = subject
		}
		eui := instructions.NewEngineUploadIntruction(flags)
		eui.Validate()
		eui.Execute()
	} else {
		util.Fail("Unknown command")

	}
}
