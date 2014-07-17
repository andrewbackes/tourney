package main

import (
	"bufio"
	// "bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// type labelwriter struct {
// 	label  string
// 	writer *bufio.Writer
// }

// func (lw labelwriter) Write(p []byte) (int, error) {
// 	s := string(p[:len(p)])
// 	n, e := lw.writer.WriteString(lw.label + s)
// 	defer lw.writer.Flush()
// 	return lw.writer.WriteString(lw.label + s)
// }

// func NewLabelWriter(basewriter io.Writer, label string) labelwriter {
// 	newwriter := bufio.NewWriter(basewriter)
// 	return labelwriter{label, newwriter}
// }

func NewNode(StdinChan chan string, command string, outHandler io.Writer, errHandler io.Writer) {
	cmd := exec.Command(command)
	StdinPipe, _ := cmd.StdinPipe()
	writer := bufio.NewWriter(StdinPipe)
	cmd.Stdout = outHandler
	cmd.Stderr = errHandler
	go cmd.Run()
	for line := range StdinChan {
		writer.WriteString(line)
		writer.Flush()
	}
}

func main() {
	StdIn := make([]chan string, 16)
	Commands := make([]string, 16)
	n := 0
	for reader := bufio.NewReader(os.Stdin); true; {
		line, _ := reader.ReadString('\n')
		address_length := strings.IndexRune(line, ' ')
		if address_length == -1 {
			continue
		}
		address, _ := strconv.Atoi(line[1:address_length])
		command := line[address_length+1 : len(line)-1]
		if line[0] == 'p' {
			if address < n {
				StdIn[address] <- command
			}
		} else if line[0] == 's' {
			Commands[n] = command
			StdIn[n] = make(chan string)
			go NewNode(StdIn[n], command[0:len(command)], os.Stdout, os.Stderr)
			fmt.Printf("Starting command `%s` on pipe %d.\n", command, n)
			n += 1
		} //else if line[0:len(line)] == "ls" {
		//	for i := 0; i < n; i += 1 {
		//		fmt.Printf("%s: %s", i, Commands[i])
		//	}
		//}
	}
}
