/*******************************************************************************

 Project: Tourney

 Module: engines
 Description: Engine struct, Protocoler interface, and UCI/WinBoard structs.

 The Engine object has a Protocoler member. Then structs corresponding to UCI
 and Winboard impliment the Protocoler interface. The engine executable itself
 is ran in a goroutine, reader and writer Engine data members read/write the
 executables stdio so other Engine methods can interact with the executable.

TODO:
	-Error checking for if engine path exists and if it opens okay
	-Engines need to take options for hashtable size, multithreading, pondering,
	opening book, and a few other bare minimums.
	-Fix the bug where >> >> >> >> ... keeps looping sometimes.

 Author(s): Andrew Backes
 Created: 7/16/2014

*******************************************************************************/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	//"strconv"
	"path/filepath"
	"strings"
	"time"
)

// helper:
type rec struct {
	timestamp time.Time
	data      string
}

/*******************************************************************************

	Protocol Specific:

*******************************************************************************/

type Protocoler interface {
	Initialize() (string, func(string) bool)
	Move(timers [2]int64, BonusTime int64, MovesToGo int64, EngineColor Color) (string, func(string) bool)
	Ping(int) (string, func(string) bool)

	NewGame(Time, Moves int64) string
	SetBoard(moveSoFar []Move, analysisSoFar []MoveAnalysis) string
	Quit() string

	ExtractMove(string) (Move, MoveAnalysis)
	RegisterEngineOptions(string, map[string]Setting)
}

/*******************************************************************************

	General Engine Functionality:

*******************************************************************************/

type Engine struct {
	//Public:
	Name string
	Path string // file location
	Spec BuildSpec
	// TODO: case sensitive issues with protocol?
	Protocol string            // = "UCI" or "WINBOARD" (TODO: Auto)
	Options  map[string]string // initialized in E.Initialize()
	MD5      string

	//Private:
	reader *bufio.Reader
	writer *bufio.Writer
	logbuf *string

	protocol         Protocoler         // should = UCI{} or WINBOARD{}
	supportedOptions map[string]Setting // decided after the engine loads and says what it supports.
}

// this struct may need to change depending on how winboard works:
type Setting struct {
	Value string
	Type  string
	Min   string
	Max   string
}

func (E *Engine) BuildRequired() bool {
	hasrepo := (E.Spec.Repo != "")
	hasbranch := (E.Spec.Branch != "")
	hasbuildfile := (E.Spec.BuildFile != "")
	return hasrepo && hasbranch && hasbuildfile
	/*
		if hasrepo && hasbranch && hasbuildfile {
			enginefile, _ := filepath.Abs( filepath.Join( E.Spec.Dir(), E.Spec.EngineFile ) )
			_, err := os.Stat(enginefile)
			built := !os.IsNotExist(err)
			return !built
		}
		return false
	*/
}

func (E *Engine) Exists() bool {
	_, err := os.Stat(E.Path)
	return !os.IsNotExist(err)
}

func (E *Engine) Equals(E2 *Engine) bool {
	return (E.Name == E2.Name) && /*(E.Path == E2.Path) &&*/ (E.Protocol == E2.Protocol)
}

func (E *Engine) ValidateEngineFile() error {
	// First decides if the file exists.
	// Compares the checksum to the md5 sum that is stored in memory.
	// If nothing has been stored, then it saves this checksum.
	// Returns true when they match or it was previously blank.
	// Primarly used when transfering a file to a worker.

	// Existence:
	if _, err := os.Stat(E.Path); os.IsNotExist(err) {
		return err
	}

	// Check sum:
	if checksum, err := GetMD5(E.Path); err != nil {
		return err
	} else {
		if E.MD5 == "" {
			// md5 has not been previously checked
			E.MD5 = checksum
		} else if E.MD5 != checksum {
			return errors.New("MD5 mismatch")
		}
	}

	return nil
}

func (E *Engine) Log(label string, record rec) {
	*E.logbuf += fmt.Sprintln("[" + record.timestamp.Format("01/02/2006 15:04:05.000") + "][" + E.Name + "][" + label + "]" + record.data)

	//DEBUG ONLY:
	//fmt.Println("[" + record.timestamp.Format("01/02/2006 15:04:05.000") + "][" + E.Name + "][" + label + "]" + record.data)
}

func (E *Engine) LogError(description string) {
	E.Log("ERROR", rec{time.Now(), description})
}

// Send a command to the engine:
func (E *Engine) Send(s string) error {
	E.Log("<-", rec{time.Now(), s})
	E.writer.WriteString(fmt.Sprintln(s)) // hopefully the line return is OS specific here.
	E.writer.Flush()
	//fmt.Print("->", fmt.Sprintln(s))
	return nil
}

// Set the engine up to be ready to think on its first move:
func (E *Engine) Start(logbuffer *string) error {
	E.logbuf = logbuffer

	// Decide which protocol to use:
	// TODO: add some autodetect code here
	if strings.ToUpper(E.Protocol) == "UCI" {
		E.protocol = &UCI{}
	} else if strings.ToUpper(E.Protocol) == "WINBOARD" {
		E.protocol = &WINBOARD{}
	}

	fullpath, _ := filepath.Abs(E.Path)
	cmd := exec.Command(fullpath)
	cmd.Dir, _ = filepath.Abs(filepath.Dir(E.Path))

	// Setup the pipes to communicate with the engine:
	StdinPipe, errIn := cmd.StdinPipe()
	if errIn != nil {
		E.LogError("Initializing Engine:" + errIn.Error())
		return errors.New("Error Initializing Engine: can not establish inward pipe.")
	}
	StdoutPipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		E.LogError("Initializing Engine:" + errOut.Error())
		return errors.New("Error Initializing Engine: can not establish outward pipe.")
	}
	E.writer, E.reader = bufio.NewWriter(StdinPipe), bufio.NewReader(StdoutPipe)

	// Start the engine:
	started := make(chan struct{})
	errChan := make(chan error)
	go func() {
		// Question: Does this force the engine to run in its own thread?
		if err := cmd.Start(); err != nil {
			errChan <- err
			return
			//return errors.New("Error executing " + E.Path + " - " + err.Error())
		}
		close(started)
	}()
	select {
	case <-started:
	case e := <-errChan:
		E.LogError("Starting Engine:" + e.Error())
		//E.Shutdown()
		return errors.New("Error starting engine:" + e.Error())
	}

	// Get the engine ready:
	if err := E.Initialize(); err != nil {
		E.LogError("Initializing Engine: " + err.Error())
		//E.Shutdown()
		//cmd.Process.Kill()
		//return err
	}

	//E.NewGame()

	// Setup up for when the engine exits:
	go func() {
		cmd.Wait()
		//TODO: add some confirmation that the engine has terminated correctly.
	}()

	return nil
}

// Send the first commands to the engine and recieves what options/features the engine supports
func (E *Engine) Initialize() error {

	s, r := E.protocol.Initialize()
	E.Send(s)
	var output string
	var err error

	output, _, err = E.Recieve(r, 2000)
	if err != nil {
		return err
	}

	// Listen to what options the engine says it supports.
	E.supportedOptions = make(map[string]Setting)
	E.protocol.RegisterEngineOptions(output, E.supportedOptions)

	E.Ping()

	return nil
}

// Recieve and process commands until a certain command is recieved
// or after the timeout (milliseconds) is achieved.
// Returns: engine output, lapsed Time, error
func (E *Engine) Recieve(EndOfRecieve func(string) bool, timeout int64) (string, time.Duration, error) {

	//var err error
	var output string //engine's output

	// Set up the Timer:
	startTime := time.Now()

	// Start recieving from the engine:
	for {
		recieved := make(chan rec, 1)
		errChan := make(chan error, 1)

		//TODO: Redesign neccessary
		go func() {
			// TODO: need a better idea here. ReadString() could hault this goroutine.
			if nextline, err := E.reader.ReadString('\n'); err == nil {
				t := time.Now()
				recieved <- rec{t, nextline}
			} else {
				errChan <- err
			}
		}()

		// Since the Timer and the reader are in goroutines, wait for:
		// (1) Something from the engine, (2) Too much Time to pass. (3) An error
		select {
		case line := <-recieved:
			// keep track of the total output from the engine:
			output += line.data // TODO: this cant possibly be fast enough

			// Take off line return bytes:
			line.data = strings.Trim(line.data, "\n") // for *nix/bsd
			line.data = strings.Trim(line.data, "\r") // for windows

			// Log this line of engine output:
			E.Log("->", line)

			// Check if the recieved command is the one we are waiting for:
			if EndOfRecieve(line.data) {
				return output, line.timestamp.Sub(startTime), nil
			}

		case <-time.After(time.Duration(timeout) * time.Millisecond):
			description := "Timed out waiting for engine to respond."
			return output, time.Now().Sub(startTime), errors.New(description)

		case e := <-errChan:
			description := "Error recieving from engine: " + e.Error()
			return output, time.Now().Sub(startTime), errors.New(description)

		}
	}
	return output, time.Now().Sub(startTime), nil //this should never occur
}

func (E *Engine) NewGame(Time, Moves int64) error {
	E.Send(E.protocol.NewGame(Time, Moves))
	return nil
}

// The engine should close itself:
func (E *Engine) Shutdown() error {
	// TODO: add confirmation that the engine has shut down correctly
	E.Send(E.protocol.Quit())
	return nil
}

// The engine should decide what move it wants to make:
func (E *Engine) Move(timers [2]int64, BonusTime int64, MovesToGo int64, EngineColor Color) (Move, MoveAnalysis, time.Duration, error) {
	s, r := E.protocol.Move(timers, BonusTime, MovesToGo, EngineColor)
	E.Send(s)
	max := timers[WHITE]
	if timers[BLACK] > max {
		max = timers[BLACK]
	}
	response, t, err := E.Recieve(r, max+1000)

	if err != nil {
		E.LogError("Requesting move: " + err.Error())
		return "", MoveAnalysis{}, t, err
	}
	// figure out what move was picked:
	mv, ma := E.protocol.ExtractMove(response)
	return mv, ma, t, nil
}

// The engine should set its internal Board to adjust for the Moves far in the game
func (E *Engine) Set(movesSoFar []Move, analysisSoFar []MoveAnalysis) error {
	s := E.protocol.SetBoard(movesSoFar, analysisSoFar)
	err := E.Send(s)
	return err
}

func (E *Engine) Ping() error {
	s, r := E.protocol.Ping(1)
	var err error
	if s != "" {
		E.Send(s)
		_, _, err = E.Recieve(r, 2000)
	}
	return err
}
