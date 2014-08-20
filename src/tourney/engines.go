/*******************************************************************************

 Project: Tourney

 Module: engines
 Description: Engine struct, Protocoler interface, and UCI/WinBoard structs.

 The Engine object has a Protocoler member. Then structs corresponding to UCI
 and Winboard impliment the Protocoler interface. The engine executable itself
 is ran in a goroutine, reader and writer Engine data members read/write the
 executables stdio so other Engine methods can interact with the executable.

TODO:
	-Error checking for it engine path exists and if it opens okay
	-WinBoard protocoler
	-Engines need to take options for hashtable size, multithreading, pondering,
	opening book, and a few other bare minimums.
	-Fix the bug where >> >> >> >> ... keeps looping sometimes.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*******************************************************************************/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*******************************************************************************

	General Engine Functionality:

*******************************************************************************/

type Engine struct {
	//Public:
	Name     string
	Path     string // file location
	Protocol string // = "UCI" or "WINBOARD"

	//Private:
	reader   *bufio.Reader
	writer   *bufio.Writer
	protocol Protocoler // should = UCI{} or WINBOARD{}
}

// Set the engine up to be ready to think on its first move:
func (E *Engine) Start() error {
	if strings.ToUpper(E.Protocol) == "UCI" {
		E.protocol = UCI{}
	} else if strings.ToUpper(E.Protocol) == "WINBOARD" {
		//E.protocol = WINBOARD{}
	}

	cmd := exec.Command(E.Path)
	// TODO: need error handling for whether or not the program launched:
	go cmd.Run()
	// TODO: wait for some sort of 'ok' to be recieved from the program that launched

	StdinPipe, errIn := cmd.StdinPipe()
	if errIn != nil {
		return errors.New("Error Initializing Engine: can not establish in pipe.")
	}
	StdoutPipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		return errors.New("Error Initializing Engine: can not establish out pipe.")
	}
	E.writer, E.reader = bufio.NewWriter(StdinPipe), bufio.NewReader(StdoutPipe)

	E.protocol.Initialize(E.reader, E.writer)
	// TODO: Set options here
	E.NewGame()

	return nil
}

func (E *Engine) NewGame() error {
	E.protocol.New(E.reader, E.writer)
	return nil
}

// The engine should close itself:
func (E *Engine) Shutdown() error {
	return E.protocol.Quit(E.writer)
}

// The engine should decide what move it wants to make:
func (E *Engine) Move(timers [2]int64, movesToGo int64) (Move, error) {
	return E.protocol.Move(E.reader, E.writer, timers, movesToGo)
}

// The engine should set its internal board to adjust for the moves far in the game
func (E *Engine) Set(movesSoFar []Move) error {
	return E.protocol.Set(E.writer, movesSoFar)
}

/*******************************************************************************

	Protocol Specific:

*******************************************************************************/

type Protocoler interface {
	// Set engine options and put it in a state to take a game position for the first time:
	Initialize(reader *bufio.Reader, writer *bufio.Writer) error
	Quit(writer *bufio.Writer) error
	Move(reader *bufio.Reader, writer *bufio.Writer, timers [2]int64, movesToGo int64) (Move, error)
	Set(writer *bufio.Writer, movesSoFar []Move) error //temporary
	New(reader *bufio.Reader, writer *bufio.Writer) error
}

type UCI struct{}
type WINBOARD struct{}

/*******************************************************************************

	UCI:

*******************************************************************************/

func (U UCI) Initialize(reader *bufio.Reader, writer *bufio.Writer) error {
	var line string

	fmt.Print("> uci\n")
	/*
		if _, err := writer.WriteString("uci\n"); err != nil {
			return errors.New("Error initializing engine. Engine not ready to accept input.")
		}
	*/
	writer.WriteString("uci\n")
	writer.Flush()

	// TODO: sometimes this loop goes infinite! probably has something to do with
	//			the time it takes the engine to load in the beginning.
	//			Does the protocol require a 1 second delay here?
	startTime := time.Now()
	for line != "uciok\n" {
		line, _ = reader.ReadString('\n')
		if line != "" {
			fmt.Print(">> ", line)
		}
		// Allow 1 second before timing out:
		if time.Now().Sub(startTime).Seconds() > 1 {
			return errors.New("Timed out. Did not recieve 'uciok' response from engine.")
		}
	}

	return nil
}

func (U UCI) New(reader *bufio.Reader, writer *bufio.Writer) error {

	fmt.Print("> ucinewgame\n")
	writer.WriteString("ucinewgame\n")
	writer.Flush()

	return nil
}

func (U UCI) Quit(writer *bufio.Writer) error {
	fmt.Print("> quit\n")
	writer.WriteString("quit\n")
	writer.Flush()
	return nil
}

func (U UCI) Move(reader *bufio.Reader, writer *bufio.Writer, timer [2]int64, movesToGo int64) (Move, error) {
	goString := "go"

	if timer[WHITE] > 0 {
		goString += " wtime " + strconv.FormatInt(timer[WHITE], 10)
	}
	if timer[BLACK] > 0 {
		goString += " btime " + strconv.FormatInt(timer[BLACK], 10)
	}
	if movesToGo > 0 {
		goString += " movestogo " + strconv.FormatInt(movesToGo, 10)
	}
	goString += "\n"

	var maxTime int64
	if timer[0] >= timer[1] {
		maxTime = timer[0]
	} else {
		maxTime = timer[1]
	}

	//fmt.Print("> " + goString)
	writer.WriteString(goString)
	writer.Flush()

	var line string
	var m Move
	startTime := time.Now()
	for strings.HasPrefix(line, "bestmove") == false {
		if int64(time.Now().Sub(startTime).Seconds()*1000) > maxTime {
			return Move{Algebraic: "none"}, errors.New("Engine timed out.")
		}
		line, _ = reader.ReadString('\n')
		m.log = append(m.log, line)
	}
	m.Algebraic = strings.Split(line, " ")[1]
	m.Algebraic = strings.TrimSuffix(m.Algebraic, "\n")

	return m, nil
}

func (U UCI) Set(writer *bufio.Writer, movesSoFar []Move) error {

	var ml []string

	for _, m := range movesSoFar {
		ml = append(ml, m.Algebraic)
	}

	var pos string
	if len(movesSoFar) > 0 {
		pos = strings.Join([]string{"position startpos moves ", strings.Join(ml, " "), "\n"}, "")
	} else {
		pos = "position startpos\n"
	}

	//fmt.Print("> ", pos)
	writer.WriteString(pos)
	writer.Flush()

	return nil
}

/*******************************************************************************

	WINBOARD:

*******************************************************************************/

/* TODO: EVERYTHING */
