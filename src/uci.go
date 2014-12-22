/*******************************************************************************

 Project: Tourney

 Module: uci
 Description: implements Protocoler from engines.go

TODO:


 Author(s): Andrew Backes
 Created: 12/15/2014

*******************************************************************************/

package main

import (
	//"bufio"
	//"errors"
	//"fmt"
	//"os"
	//"os/exec"
	"strconv"
	"strings"
	//"time"
)

/*******************************************************************************

	UCI:

*******************************************************************************/

type UCI struct{}

func (U UCI) Initialize() (string, func(string) bool) {
	// (command to send),(command to recieve)
	recieve := "uciok"
	return "uci", func(parse string) bool {
		return len(parse) >= len(recieve) && parse[:len(recieve)] == recieve
	}
}

func (U UCI) Move(Timer [2]int64, MovesToGo int64, EngineColor Color) (string, func(parse string) bool) {
	goString := "go"

	if Timer[WHITE] > 0 {
		goString += " wtime " + strconv.FormatInt(Timer[WHITE], 10)
	}
	if Timer[BLACK] > 0 {
		goString += " btime " + strconv.FormatInt(Timer[BLACK], 10)
	}
	if MovesToGo > 0 {
		goString += " movestogo " + strconv.FormatInt(MovesToGo, 10)
	}
	recieve := "bestmove "
	return goString, func(parse string) bool {
		return len(parse) >= len(recieve) && parse[:len(recieve)] == recieve
	}
}

func (U UCI) Ping(N int) (string, func(parse string) bool) {
	recieve := "readyok"
	return "isready", func(parse string) bool {
		return len(parse) >= len(recieve) && parse[:len(recieve)] == recieve
	}
}

func (U UCI) NewGame(Time, Moves int64) string {
	return "ucinewgame"
}

func (U UCI) SetBoard(movesSoFar []Move) string {
	var ml []string

	for _, m := range movesSoFar {
		ml = append(ml, m.Algebraic)
	}

	var pos string
	if len(movesSoFar) > 0 {
		pos = "position startpos moves " + strings.Join(ml, " ")
	} else {
		pos = "position startpos"
	}

	return pos
}

func (U UCI) Quit() string {
	return "quit"
}

func (U UCI) RegisterEngineOptions(output string, options map[string]Setting) {
	// TODO: some engines have two word setting keys
	// 		 ex: senpai: option name Log File type check default false
	// TODO: case sensitive

	if output == "" {
		return
	}

	output = strings.Replace(output, "\r", "", -1) // remove the \r after the \n\r
	lines := strings.Split(output, "\n")
	for i, _ := range lines {
		newSettingLabel := ""
		newSetting := Setting{}
		words := strings.Split(lines[i], " ")
		// double check that this line has option information on it:
		if strings.ToLower(words[0]) != "option" {
			continue
		}
		// Process the option information:
		for j := 0; j < len(words)-1; j++ {
			switch strings.ToLower(words[j]) {
			case "name":
				newSettingLabel = words[j+1]
			case "type":
				newSetting.Type = words[j+1]
			case "default":
				newSetting.Value = words[j+1]
			case "min":
				newSetting.Min = words[j+1]
			case "max":
				newSetting.Max = words[j+1]
			}
		}
		options[newSettingLabel] = newSetting
	}
}

func parseInfo(line string) EvaluationData {
	// ex: info depth 120 seldepth 2 time 10 nodes 6792 nps 679200 multipv 1 score mate 1 pv c5d7
	// ex: info multipv 1 depth 4 score mate -5 time 4 nodes 7159 hashfull 321 pv f5f4 d7d4 f4g3 d4e3 g3g4 e3f3 g4g5 c5e6 g5h6 f8h8
	// ex: info multipv 1 depth 9 score cp -1261 upperbound time 7 nodes 11461 hashfull 286 pv h7g7 d2h6 g7f7 h6f8 f7e6 f8c8 e6e5 c8c5 e5e6 c5d6 e6f7 h8h7 f7g8 d6e7 h4e1 g1h2 e1h4 g2h3 h4f4 h2g2

	if !strings.HasPrefix(line, "info") {
		return EvaluationData{}
	}

	words := strings.Split(line, " ")

	var d, sd, n, t, s int
	var l, u bool
	var pv string

	for i := 1; i < len(words)-1; i++ {
		switch words[i] {
		case "depth":
			d, _ = strconv.Atoi(words[i+1])
		case "seldepth":
			sd, _ = strconv.Atoi(words[i+1])
		case "nodes":
			n, _ = strconv.Atoi(words[i+1])
		case "time":
			t, _ = strconv.Atoi(words[i+1])
		case "score":
			if i+2 < len(words) {
				if words[i+1] == "cp" {
					s, _ = strconv.Atoi(words[i+2])
				} else if words[i+1] == "mate" {
					s, _ = strconv.Atoi(words[i+2])
					s = MateIn(s)
				}
			}
			if i+3 < len(words) {
				l = (words[i+3] == "lowerbound")
				u = (words[i+3] == "upperbound")
			}
		case "pv":
			for j := i + 1; j < len(words); j++ {
				if isMove(words[j]) {
					pv += words[j] + " "
				} else {
					break
				}
			}
		}
	}

	return EvaluationData{Depth: d, Seldepth: sd, Nodes: n, Score: s, Upperbound: u, Lowerbound: l, Time: t, Pv: strings.Trim(pv, " ")}
}

func (U UCI) ExtractMove(output string) Move {

	output = strings.Replace(output, "\r", "", -1)
	lines := strings.Split(output, "\n")
	mv := Move{}
	for _, line := range lines {
		eval := parseInfo(line)
		if eval.Depth != 0 {
			mv.Evaluation = append(mv.Evaluation, eval)
		}
		if strings.HasPrefix(line, "bestmove") {
			words := strings.Fields(line)
			if len(words) >= 2 {
				mv.Algebraic = words[1]
			}
			if len(words) >= 4 && words[2] == "ponder" {
				mv.Ponder = words[3]
			}
		}
	}
	return mv

	/*
		// TODO: REFACTOR: this replace also happens in Engine.Recieve()
		output = strings.Replace(output, "\n\r", " ", -1)
		output = strings.Replace(output, "\n", " ", -1)

		words := strings.Split(output, " ")

		// Helper functions:
		LastNValuesOf := func(key string, N int) string {
			for i := len(words) - 1; i >= 0; i-- {
				if words[i] == key {
					if i+N <= len(words)-1 {
						return strings.Join(words[i+1:i+N+1], " ")
					}
				}
			}
			return ""
		}
		LastValueOf := func(key string) string {
			//returns the word after the word given as an arg
			return LastNValuesOf(key, 1)
		}

		// ***

		keys := []string{"depth", "time", "nodes"}
		values := [4]int{0, 0, 0, 0}
		for i, key := range keys {
			temp := LastValueOf(key)
			if isNumber(temp) {
				values[i], _ = strconv.Atoi(temp)
			}
		}
		skey := LastValueOf("score")
		var sval int
		if skey == "cp" {
			sval, _ = strconv.Atoi(LastValueOf(skey))
		} else if skey == "mate" {
			sval, _ = strconv.Atoi(LastValueOf(skey))
			sval = MateIn(sval)
		}

		return (Move{
			Algebraic: LastValueOf("bestmove"),
			Depth:     values[0],
			Time:      values[1],
			Nodes:     values[2],
			Score:     sval,
			Pv:        LastNValuesOf("pv", values[0]),
		})
	*/
}
