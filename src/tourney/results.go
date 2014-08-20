/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/11/2014

 Module: Results
 Description: Functions that have to do with tourney records. Gathering the
 			  records, formatting, etc.

*******************************************************************************/

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Record struct {
	Player     Engine
	Opponent   Engine
	Wins       int
	Losses     int
	Draws      int
	Incomplete int
	Order      string // w-l-d string. example: 10=11=01
}

// Returns the scores for each different matchup in the tourney:
func MatchupResults(T *Tourney) []Record {
	var r []Record

	// helper function:
	indexOf := func(e Engine, o Engine) int {
		for i, _ := range r {
			if r[i].Player.Name == e.Name && r[i].Opponent.Name == o.Name {
				return i
			}
		}
		r = append(r, Record{Player: e, Opponent: o})
		return len(r) - 1
	}
	// workhorse:
	for i, _ := range T.GameList {
		for color := WHITE; color <= BLACK; color++ {
			ind := indexOf(T.GameList[i].Player[color], T.GameList[i].Player[[]Color{BLACK, WHITE}[color]])
			if T.GameList[i].Completed {
				if winner := T.GameList[i].Result; winner == DRAW {
					r[ind].Draws++
					r[ind].Order += "="
				} else if winner == color {
					r[ind].Wins++
					r[ind].Order += "1"
				} else {
					r[ind].Losses++
					r[ind].Order += "0"
				}
			} else {
				r[ind].Incomplete++
				r[ind].Order += "?"
			}
		}
	}

	// helper:
	score := func(rec Record) int {
		if (rec.Wins + rec.Draws + rec.Losses) == 0 {
			return 0
		}
		return (10000*rec.Wins + 5000*rec.Draws) / (rec.Wins + rec.Draws + rec.Losses)
	}
	// Sort by Player.Name and then by highest %
	// Group records according to their player names:
	for i := 0; i < len(r)-1; i++ {
		pivot := i + 1
		if r[i].Player.Name == r[pivot].Player.Name {
			continue
		}
		for j := i + 2; j < len(r)-1; j++ {
			if r[i].Player.Name == r[j].Player.Name {
				placeHolder := r[j]
				r[j] = r[pivot]
				r[pivot] = placeHolder
				break
			}
		}
	}
	// within the grouping by name, sort by score:
	var begin, end int
	for {
		for end = begin; end < len(r)-1; end++ {
			if r[begin].Player.Name != r[end+1].Player.Name {
				break
			}
		}
		//bubble sort within this group:
		for e := int(end); e >= begin; e-- {
			for i := begin; i <= e-1; i++ {
				if score(r[i]) < score(r[i+1]) {
					//swap
					placeholder := r[i+1]
					r[i+1] = r[i]
					r[i] = placeholder
				}
			}
		}
		//skip past the current grouping:
		begin = end + 1
		if begin > len(r) {
			break
		}
	}
	return r
}

// Returns the results for each individual engine in the tourney:
func EngineResults(T *Tourney) []Record {
	// TODO: this function has a lot of overlap with MatchupResults(). Should refactor.
	var r []Record

	// helper function:
	indexOf := func(e Engine) int {
		for i, _ := range r {
			if r[i].Player.Name == e.Name {
				return i
			}
		}
		r = append(r, Record{Player: e})
		return len(r) - 1
	}
	// workhorse:
	for i, _ := range T.GameList {
		for color := WHITE; color <= BLACK; color++ {
			ind := indexOf(T.GameList[i].Player[color])
			if T.GameList[i].Completed {
				if winner := T.GameList[i].Result; winner == DRAW {
					r[ind].Draws++
					r[ind].Order += "="
				} else if winner == color {
					r[ind].Wins++
					r[ind].Order += "1"
				} else {
					r[ind].Losses++
					r[ind].Order += "0"
				}
			} else {
				r[ind].Incomplete++
				r[ind].Order += "?"
			}
		}
	}

	// Sort by highest %. I was lazy and just did a bubble sort:
	score := func(rec Record) int {
		if (rec.Wins + rec.Draws + rec.Losses) == 0 {
			return 0
		}
		return (10000*rec.Wins + 5000*rec.Draws) / (rec.Wins + rec.Draws + rec.Losses)
	}
	for end := int(len(r) - 1); end >= 0; end-- {
		for i := 0; i <= end-1; i++ {
			if score(r[i]) < score(r[i+1]) {
				//swap
				placeholder := r[i+1]
				r[i+1] = r[i]
				r[i] = placeholder
			}
		}
	}

	return r
}

// Takes a record and spits out a string of what that record contains:
func FormatRecord(record Record) string {
	var str, matchup string
	// Engine Names:
	if record.Opponent.Name == "" {
		matchup = record.Player.Name
	} else {
		matchup = record.Player.Name + " - " + record.Opponent.Name
	}
	str += fmt.Sprint(matchup, strings.Repeat(" ", 40-len(matchup)), ":   ")
	// W-L-D :
	str += fmt.Sprint(record.Wins, "-", record.Losses, "-", record.Draws, "\t")
	// Point score:
	score := float64(record.Wins) + 0.5*float64(record.Draws)
	possible := float64(record.Wins + record.Losses + record.Draws)
	// As fraction:
	str += fmt.Sprint(score, "/", possible, "\t")
	// As percentage:
	if possible > 0 {
		str += fmt.Sprintf("%.2f", 100*(score/possible))
		str += "%"
	} else {
		str += "00.00%"
	}
	// win-loss-draw single line chart:
	/*
		if len(record.Order) < 36 {
			str += strings.Repeat(" ", 43)
		} else if len(record.Order) <= 68 {
			str += strings.Repeat(" ", 70-len(record.Order))
		}
	*/
	//if l := 80 - len(str); l > 0 {
	str += strings.Repeat(" ", 11)
	//str += strconv.Itoa(l)
	//}
	str += "(" + record.Order + ")\n"
	return str
}

// Creates a report summarizing the tourney results.
// Scores for each matchup in the tourney and also each engine's overall score.
func SummarizeResults(T *Tourney) string {

	matchups := MatchupResults(T)
	engines := EngineResults(T)
	matchupSummary := strings.Repeat("=", 80) + "\n   Results by Matchup:\n" + strings.Repeat("=", 80) + "\n"
	for i, _ := range matchups {
		if i > 0 && matchups[i].Player.Name != matchups[i-1].Player.Name {
			matchupSummary += strings.Repeat("-", 80) + "\n"
		}
		matchupSummary += FormatRecord(matchups[i])
	}
	eventSummary := strings.Repeat("=", 80) + "\n   Event Summary:\n" + strings.Repeat("=", 80) + "\n"
	for _, record := range engines {
		eventSummary += FormatRecord(record)
	}
	// count completed games:
	completed := 0
	for _, g := range T.GameList {
		if g.Completed {
			completed++
		}
	}
	eventSummary += strings.Repeat("-", 80) + "\nGames played: " + strconv.Itoa(completed) + "/" + strconv.Itoa(len(T.GameList)) + "\n"
	return matchupSummary + "\n" + eventSummary + "\n"
}

func SummarizeGames(T *Tourney) string {
	// Event, Round, Site, Date, White, Black, Result, Details
	summary := strings.Repeat("=", 80) + "\n   Game History:\n" + strings.Repeat("=", 80) + "\n"
	for _, g := range T.GameList {
		summary += g.Event + ", " +
			strconv.Itoa(g.Round) + ", " +
			g.Site + ", " +
			g.Date + ", " +
			g.Player[WHITE].Name + ", " +
			g.Player[BLACK].Name + ", "
		if g.Completed {
			summary += []string{"1-0", "0-1", "1/2-1/2"}[g.Result] + ", "
		} else {
			summary += "*, "
		}
		summary += g.ResultDetail + "\n"
	}
	return summary
}
