/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 8/11/2014

 Module: Results
 Description: Functions that have to do with tourney records. Gathering the
 			  records, formatting, etc.

TODO:
	-Refactor. Too much similarity between functions. Combine.
	-Rework the result rollup in such a way that the html/template's are more
	 user friendly.
	-Use text/templates to display on screen results as well as save to files.
	-rename this to standings.

*******************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type TourneyStandings struct {
	//			   [player]   [opponent]
	records     map[string]map[string]*PlayerRecord
	orderedKeys map[string][]string
}

type PlayerRecord struct {
	Name                  string
	Wins                  int64
	Losses                int64
	Draws                 int64
	Incomplete            int64
	Graph                 []byte
	WinningConditionCount map[string]int64
	DrawConditionCount    map[string]int64
	LossingConditionCount map[string]int64
}

var DRAW_CONDITIONS = []string{
	STALEMATE, FIFTY_MOVE, THREE_FOLD, INSUFFICIENT_MATERIAL,
}
var WIN_LOSS_CONDITIONS = []string{
	CHECKMATE, TIMED_OUT, STOPPED_RESPONDING, ILLEGAL_MOVE,
}

func (P PlayerRecord) TrimmedName() string {
	invalid := []string{"~", "!", "@", "$", "%", "^", "&", "*", "(", ")", "+", "=", ",", "."}
	trimmed := P.Name
	for _, v := range invalid {
		strings.Replace(trimmed, v, "", -1)
	}
	return trimmed
}

func (P PlayerRecord) Score() float32 {
	return float32(P.Wins) + (0.5 * float32(P.Draws))
}

func (P PlayerRecord) Rate() float32 {
	n := 10000*P.Wins + 5000*P.Draws
	d := P.TotalGames()
	if d == 0 {
		return 0
	}
	return float32(n/d) / float32(100)
}

func (P PlayerRecord) TotalGames() int64 {
	return P.Wins + P.Losses + P.Draws
}

func NewTourneyStandings(T *Tourney) *TourneyStandings {
	TourneyResults := TourneyStandings{}
	TourneyResults.records = make(map[string]map[string]*PlayerRecord)
	TourneyResults.records["All"] = make(map[string]*PlayerRecord)
	for j, _ := range T.Engines {
		TourneyResults.records[T.Engines[j].Name] = make(map[string]*PlayerRecord)
	}
	return &TourneyResults
}

//
// Goes through the Tourney's game list and generates all of the records.
// This function should not really be needed since after each game
// the records should be updated.
//
func GenerateGameRecords(T *Tourney, drawGraph bool) *TourneyStandings {
	TourneyResults := NewTourneyStandings(T)
	for i, _ := range T.GameList {
		TourneyResults.AddOrUpdateGame(&T.GameList[i], drawGraph, false)
	}
	//TourneyResults.SortKeys()
	return TourneyResults
}

//
// Takes the information from a Game struct and adds it to the standings
// for the tournament. Adjusts/Sorts the rankings accordingly.
//
// The update flag indicates what to do about incomplete games. If the game doesn't
// already holds a record then incompletes need to be counted. If it's not a new
// game then completed games need to subtract from the incomplete games.
//
// The updateGraph flag controls whether or not the win/loss/draw linear
// graph will be used. For very large tournaments it is best not to use
// the graph.
//
func (R *TourneyStandings) AddOrUpdateGame(G *Game, updateGraph bool, update bool) {
	w, b := G.Player[0].Name, G.Player[1].Name

	rec := []PlayerRecord{PlayerRecord{Name: b}, PlayerRecord{Name: w}, PlayerRecord{Name: w}, PlayerRecord{Name: b}}
	rec_pntr := []*PlayerRecord{R.records[w][b], R.records[b][w], R.records["All"][w], R.records["All"][b]}
	player_keys := []string{w, b, "All", "All"}
	opponent_keys := []string{b, w, w, b}

	for i, _ := range rec_pntr {
		if rec_pntr[i] != nil {
			rec[i] = *rec_pntr[i]
		} else {
			if R.orderedKeys == nil {
				R.orderedKeys = make(map[string][]string)
			}
			R.orderedKeys[player_keys[i]] = append(R.orderedKeys[player_keys[i]], opponent_keys[i])
		}
	}
	if G.Completed == false && !update {
		for i, _ := range rec {
			rec[i].Incomplete++
		}
	} else {
		if update {
			for i, _ := range rec {
				rec[i].Incomplete--
			}
		}
		if G.Result != DRAW {
			rec[G.Result].Wins++
			rec[G.Result].WinningConditionCount = IncrementMap(rec[G.Result].WinningConditionCount, G.EndingCondition())
			rec[G.Result+2].Wins++
			rec[G.Result+2].WinningConditionCount = IncrementMap(rec[G.Result+2].WinningConditionCount, G.EndingCondition())
			rec[1-G.Result].Losses++
			rec[1-G.Result].LossingConditionCount = IncrementMap(rec[1-G.Result].LossingConditionCount, G.EndingCondition())
			rec[(1-G.Result)+2].Losses++
			rec[(1-G.Result)+2].LossingConditionCount = IncrementMap(rec[(1-G.Result)+2].LossingConditionCount, G.EndingCondition())
			if updateGraph {
				rec[G.Result].Graph = append(rec[G.Result].Graph, '1')
				rec[G.Result+2].Graph = append(rec[G.Result+2].Graph, '1')
				rec[1-G.Result].Graph = append(rec[1-G.Result].Graph, '0')
				rec[(1-G.Result)+2].Graph = append(rec[(1-G.Result)+2].Graph, '0')
			}
		} else {
			for i, _ := range rec {
				rec[i].Draws++
				rec[i].DrawConditionCount = IncrementMap(rec[i].DrawConditionCount, G.EndingCondition())
				if updateGraph {
					rec[i].Graph = append(rec[i].Graph, '=')
				}
			}
		}
	}

	R.records[w][b], R.records[b][w] = &rec[0], &rec[1]
	R.records["All"][w], R.records["All"][b] = &rec[2], &rec[3]

	// Sort:
	R.SortRecordsForPlayer(w)
	R.SortRecordsForPlayer(b)
	R.SortRecordsForPlayer("All")

}

//
// Provides a list of the overall rankings of the players.
//
func (S *TourneyStandings) OverallStandings() []*PlayerRecord {
	return S.MatchupStandings("All")
}

func (S *TourneyStandings) OverallStandingsFor(player string) *PlayerRecord {
	return S.records["All"][player]
}

//
// Returns the records of all the opponents against the given player.
// The records should be sorted before hand.
//
func (S *TourneyStandings) MatchupStandings(player string) []*PlayerRecord {
	// create the slice to return:
	records := make([]*PlayerRecord, len(S.orderedKeys[player]))
	// loop through gettings the records in order:
	for i, key := range S.orderedKeys[player] {
		records[i] = S.records[player][key]
	}
	return records
}

//
// For each player, all of the opponent's standings are returned. Players
// are ordered by rank.
//
func (S *TourneyStandings) AllMatchupStandings() []*PlayerRecord {
	var records []*PlayerRecord
	for _, player := range S.Players() {
		records = append(records, S.MatchupStandings(player)...)
	}
	return records
}

//
// Provides a list of players that have standings in order of ranking
//
func (S *TourneyStandings) Players() []string {
	A := S.OverallStandings()
	r := make([]string, len(A))
	for i, _ := range A {
		r[i] = A[i].Name
	}
	return r
}

//
// Spits the player rankings out in the format given by the template.
//
func (R *TourneyStandings) RenderTemplate(filename string) string {
	file := filepath.Join(Settings.TemplateDirectory, filename)
	tmpl, err := template.ParseFiles(file)

	if err != nil {
		fmt.Println(err)
		//io.WriteString(w, fmt.Sprint("Error opening '", file, "' - ", err))
		return ""
	}
	var w bytes.Buffer
	err = tmpl.Execute(&w, R)
	if err != nil {
		fmt.Println(err)
		//io.WriteString(w, fmt.Sprint("Error executing parse on '", file, "' - ", err))
		return ""
	}
	return w.String()
}

func (R *TourneyStandings) PrintStandings() {
	fmt.Println(R.RenderTemplate("standings.txt"))
}

/*******************************************************************************

	Sorting:

*******************************************************************************/

type RecordSorter struct {
	Keys    []string
	Records map[string]*PlayerRecord
}

func (S RecordSorter) Len() int {
	return len(S.Keys)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (S RecordSorter) Less(i, j int) bool {
	a := *S.Records[S.Keys[i]]
	b := *S.Records[S.Keys[j]]
	return a.Score() > b.Score()
}

// Swap swaps the elements with indexes i and j.
func (S RecordSorter) Swap(i, j int) {
	S.Keys[i], S.Keys[j] = S.Keys[j], S.Keys[i]
}

//
// Sorts all of the records in the TourneyStandings at a matchup level.
// Highest score to lowest.
//
func (R *TourneyStandings) SortAllKeys() {
	R.orderedKeys = make(map[string][]string)
	// Populate the list of Keys:
	for player, record := range R.records {
		for opponent, _ := range record {
			R.orderedKeys[player] = append(R.orderedKeys[player], opponent)
		}
		// Sort the list of Keys based on score:
		data := RecordSorter{
			Keys:    R.orderedKeys[player],
			Records: R.records[player],
		}
		sort.Sort(data)
	}
}

//
// Sorts the opponents records of the specified played.
// Arranges highest score to lowest.
//
func (R *TourneyStandings) SortRecordsForPlayer(player string) {
	data := RecordSorter{
		Keys:    R.orderedKeys[player],
		Records: R.records[player],
	}
	sort.Sort(data)
}
