/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/8/2014

 Module: PGN
 Description: PGN tools

*******************************************************************************/

package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Turns Game structs into PGN
func EncodePGN(G *Game) string {
	var pgn string
	tags := [][]string{
		{"Event", G.Event},
		{"Site", G.Site},
		{"Date", G.Date},
		{"Round", strconv.Itoa(G.Round)},
		{"White", G.Player[WHITE].Name},
		{"Black", G.Player[BLACK].Name}}
	if G.Completed {
		tags = append(tags, []string{"Result", []string{"1-0", "0-1", "1/2-1/2"}[G.Result]})
	} else {
		tags = append(tags, []string{"Result", "*"})
	}
	for _, t := range tags {
		pgn += "[" + t[0] + " \"" + t[1] + "\"]\n"
	}
	pgn += "\n"
	for j, _ := range G.MoveList {
		if j%2 == 0 {
			pgn += strconv.Itoa((j/2)+1) + ". "
		}
		pgn += G.MoveList[j].Algebraic + " "
	}
	pgn += tags[6][1] + "\n\n"

	return pgn
}

// Turns PGN into Game Structs
func DecodePGN(pgn string) []Game {
	var gamelist []Game

	game := strings.Split(pgn, "[Event")
	for i, _ := range game {
		var G Game
		if game[i] != "" {
			game[i] = "[Event" + game[i]
		}
		lines := strings.Split(game[i], "\n")
		for _, l := range lines {
			if l == "" {
				continue
			}
			if strings.HasPrefix(l, "[") {
				// its a tag
				t := strings.SplitN(l, " ", 2)
				key := t[0][1:]
				value := t[1][1 : len(t[1])-2]
				switch key {
				case "Event":
					G.Event = value
				case "Site":
					G.Site = value
				case "Date":
					G.Date = value
				case "Round":
					G.Round, _ = strconv.Atoi(value)
				case "White":
					G.Player[WHITE].Name = value
				case "Black":
					G.Player[BLACK].Name = value
				case "Result":
					G.Result = map[string]Color{"1-0": WHITE, "0-1": BLACK, "1/2-1/2": DRAW, "*": NEITHER}[value]
				}
			} else {
				// its a move
				moves := strings.Split(l, " ")
				for _, m := range moves {
					if strings.HasPrefix(m, "{") {
						continue
					}
					SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x=]?)([BKNPQR]|[a-h][1-8])([+#!?]?)([+#!?]?)"
					if ok, _ := regexp.MatchString(SAN, m); ok || m == "O-O" || m == "O-O-O" {
						G.MoveList = append(G.MoveList, Move{Algebraic: m})
					}
				}
			}
		}
		if (G.Event != "") && (len(G.MoveList) > 0) {
			gamelist = append(gamelist, G)
		}
	}

	return gamelist
}
