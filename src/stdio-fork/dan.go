package main

import (
	"bufio"
	// "bytes"
	"fmt"
	// "io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// type LabelWriter struct {
// 	label  string
// 	Writer *bufio.Writer
// }

// func (lw LabelWriter) Write(p []byte) (int, error) {
// 	s := string(p[:len(p)])
// 	defer lw.Writer.Flush()
// 	return lw.Writer.WriteString(lw.label + s)
// }

// func NewLabelWriter(baseWriter io.Writer, label string) LabelWriter {
// 	newWriter := bufio.NewWriter(baseWriter)
// 	return LabelWriter{label, newWriter}
// }

func SpawnChild(StdinChan chan string, command string, address int) {
	cmd := exec.Command(command)
	StdinPipe, _ := cmd.StdinPipe()
	writer := bufio.NewWriter(StdinPipe)
	cmd.Stdout = os.Stdout //NewLabelWriter(os.Stdout, fmt.Sprintf("%d : ", address))
	cmd.Stderr = os.Stderr //NewLabelWriter(os.Stderr, fmt.Sprintf("%d : ", address))
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
			go SpawnChild(StdIn[n], command[0:len(command)], n)
			fmt.Printf("Starting command `%s` on pipe %d.\n", command, n)
			n += 1
		} //else if line[0:len(line)] == "ls" {
		//	for i := 0; i < n; i += 1 {
		//		fmt.Printf("%s: %s", i, Commands[i])
		//	}
		//}
	}
}
