package main

import (
	"fmt"
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

func uploadEngine(f, n, v, o string) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		fmt.Println(f + " does not exist")
		os.Exit(1)
	}
	targetURL := getAPIURL() + "/engineFiles/" + n + "/" + v + "/" + o
	err := util.PostFile(f, targetURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Uploaded", f)
}

func registerEngine(n, v, o string) {
	fmt.Println("Registering Engine\nName:    ", n, "\nVersion: ", v, "\nOS:      ", o)
}

func main() {
	verb, noun, subject := command()
	if verb == "add" && noun == "engine" && subject != "" {
		flags := flags()
		n, ok1 := flags["name"]
		v, ok2 := flags["version"]
		o, ok3 := flags["os"]
		if ok1 && ok2 && ok3 {
			uploadEngine(subject, n, v, o)
			registerEngine(n, v, o)
		} else {
			fmt.Println("'add engine' command requires --name, --version, and --os")
		}
	} else {
		fmt.Println("Unknown command")
	}
}

func getAPIURL() string {
	if os.Getenv("API_URL") != "" {
		return os.Getenv("API_URL")
	}
	return "http://api.tourney.aback.es:9090/api/v2"
}
