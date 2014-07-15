package main

import (
	"bufio"
	"os"
	"os/exec"
)

func child_loop(StdinChan chan string, command string) {
	cmd := exec.Command(command)
	StdinPipe, _ := cmd.StdinPipe()
	writer := bufio.NewWriter(StdinPipe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go cmd.Run()
	for line := range StdinChan {
		writer.WriteString(line)
		writer.Flush()
	}
}

func spawn_child_proc(command string) chan string {
	StrChan := make(chan string)
	go child_loop(StrChan, command)
	return StrChan
}

func main() {
	StdIn := spawn_child_proc("../bin/child")
	StdInAlt := spawn_child_proc("../bin/child2")
	for reader := bufio.NewReader(os.Stdin); true; {
		line, _ := reader.ReadString('\n')
		if line[0] == '1' {
			StdIn <- line[1:len(line)]
		} else if line[0] == '2' {
			StdInAlt <- line[1:len(line)]

		}
		if line == "exit" {
			close(StdIn)
			close(StdInAlt)
			break
		}
	}
}
