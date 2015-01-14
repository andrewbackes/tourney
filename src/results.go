/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/11/2014

 Module: Results
 Description: Functions that have to do with tourney records. Gathering the
 			  records, formatting, etc.

TODO:
	-Refactor. Too much similarity between functions. Combine.
	-Rework the result rollup in such a way that the html/template's are more
	 user friendly.
	-Use text/templates to display on screen results as well as save to files.

*******************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
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

type RecordRollup struct {
	//TODO: need a high level summary
	Info           *Tourney
	EngineRecords  []Record
	MatchupRecords []Record
}

func NewRecordRollup(T *Tourney) *RecordRollup {
	r := &RecordRollup{}
	r.Info = T
	r.EngineRecords = EngineResults(T)
	r.MatchupRecords = MatchupResults(T)
	return r
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
	str += "(" + record.Order + ")" + fmt.Sprintln()
	return str
}

// Creates a report summarizing the tourney results.
// Scores for each matchup in the tourney and also each engine's overall score.
func SummarizeResults(T *Tourney) string {

	matchups := MatchupResults(T)
	engines := EngineResults(T)
	matchupSummary := strings.Repeat("=", 80) + fmt.Sprintln() +
		"   Results by Matchup:" + fmt.Sprintln() +
		strings.Repeat("=", 80) + fmt.Sprintln()

	for i, _ := range matchups {
		if i > 0 && matchups[i].Player.Name != matchups[i-1].Player.Name {
			matchupSummary += strings.Repeat("-", 80) + fmt.Sprintln()
		}
		matchupSummary += FormatRecord(matchups[i])
	}
	eventSummary := strings.Repeat("=", 80) + fmt.Sprintln() +
		"   Event Summary:" + fmt.Sprintln() +
		strings.Repeat("=", 80) + fmt.Sprintln()

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
	eventSummary += strings.Repeat("-", 80) + fmt.Sprintln() +
		"Games played: " + strconv.Itoa(completed) + "/" + strconv.Itoa(len(T.GameList)) + fmt.Sprintln()
	return matchupSummary + fmt.Sprintln() + eventSummary + fmt.Sprintln()
}

func SummarizeGames(T *Tourney) string {
	// Event, Round, Site, Date, White, Black, Result, Details
	summary := strings.Repeat("=", 80) + fmt.Sprintln() +
		"   Game History:" + fmt.Sprintln() +
		strings.Repeat("=", 80) + fmt.Sprintln()
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
		summary += g.ResultDetail + fmt.Sprintln()
	}
	return summary
}

/*******************************************************************************

 REFACTOR:

*******************************************************************************/

type PlayerRecord struct {
	Wins       int64
	Losses     int64
	Draws      int64
	Incomplete int64
	Graph      []byte
}

func (P PlayerRecord) Score() float32 {
	return float32(P.Wins) + (0.5 * float32(P.Draws))
}

func (P PlayerRecord) Rate() float32 {
	n := 10000*P.Wins + 5000*P.Draws
	d := P.TotalGames()
	return float32(n/d) / float32(100)
}

func (P PlayerRecord) TotalGames() int64 {
	return P.Wins + P.Losses + P.Draws
}

type Result struct {
	Records     map[string]map[string]PlayerRecord
	OrderedKeys map[string][]string
}

func GenerateResults(T *Tourney, drawGraph bool) *Result {
	TourneyResults := Result{}
	TourneyResults.Records = make(map[string]map[string]PlayerRecord)
	TourneyResults.OrderedKeys = make(map[string][]string)
	TourneyResults.Records["All"] = make(map[string]PlayerRecord)
	for j, _ := range T.Engines {
		TourneyResults.Records[T.Engines[j].Name] = make(map[string]PlayerRecord)
	}
	for i, _ := range T.GameList {
		TourneyResults.Update(&T.GameList[i], drawGraph)
		//UpdateResultsFromGame(&TourneyResults, &T.GameList[i], drawGraph)
	}
	TourneyResults.SortKeys()
	return &TourneyResults
}
func (R *Result) Update(G *Game, updateGraph bool) {
	//func UpdateResultsFromGame(R *Result, G *Game, updateGraph bool) {
	w := G.Player[0].Name
	b := G.Player[1].Name
	rec := []PlayerRecord{
		R.Records[w][b],
		R.Records[b][w],
		R.Records["All"][w],
		R.Records["All"][b],
	}
	if G.Completed == false {
		for i, _ := range rec {
			rec[i].Incomplete++
		}
	} else if G.Result != DRAW {
		rec[G.Result].Wins++
		rec[G.Result+2].Wins++
		rec[1-G.Result].Losses++
		rec[(1-G.Result)+2].Losses++
		if updateGraph {
			rec[G.Result].Graph = append(rec[G.Result].Graph, '1')
			rec[G.Result+2].Graph = append(rec[G.Result+2].Graph, '1')
			rec[1-G.Result].Graph = append(rec[1-G.Result].Graph, '0')
			rec[(1-G.Result)+2].Graph = append(rec[(1-G.Result)+2].Graph, '0')
		}
	} else {
		for i, _ := range rec {
			rec[i].Draws++
			if updateGraph {
				rec[i].Graph = append(rec[i].Graph, '=')
			}
		}
	}
	R.Records[w][b] = rec[0]
	R.Records[b][w] = rec[1]
	R.Records["All"][w] = rec[2]
	R.Records["All"][b] = rec[3]
}

func (R *Result) RenderTemplate() string {
	file := filepath.Join(Settings.TemplateDirectory, "results.txt")
	tmpl, err := template.ParseFiles(file)

	if err != nil {
		fmt.Println(err)
		//io.WriteString(w, fmt.Sprint("Error opening '", file, "' - ", err))
		return ""
	}
	//tmpl.Funcs(template.FuncMap{"isMatchup": func(s string) bool { return s != "All" }})
	//tmpl = tmpl.Funcs(template.FuncMap{"isMatchup": func(s string) bool { return s != "All" }})
	var w bytes.Buffer
	err = tmpl.Funcs(template.FuncMap{"isMatchup": func(s string) bool { return s != "All" }}).Execute(&w, R)
	if err != nil {
		fmt.Println(err)
		//io.WriteString(w, fmt.Sprint("Error executing parse on '", file, "' - ", err))
		return ""
	}
	return w.String()
}

//
// This function requires the keys are sorted before calling.
//
/*
func (R *Result) EngineRecords() []PlayerRecord {
	var records []PlayerRecord
	keys := R.OrderedKeys["All"]
	for i, _ := range keys {
		records = append(records, R.Records["All"][key])
	}
	return records
}

func (R *Result) MatchupRecords() []PlayerRecord {
	var records []PlayerRecord

	for players, recs := range R.Records {
		if player == "All" {
			continue
		}
		for opponents, rec := range R.Records[players] {
			records = append()
		}
	}
}
*/

/*******************************************************************************

	Sorting:

*******************************************************************************/

type RecordSorter struct {
	Keys    []string
	Records map[string]PlayerRecord
}

func (S RecordSorter) Len() int {
	return len(S.Keys)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (S RecordSorter) Less(i, j int) bool {
	a := S.Records[S.Keys[i]]
	b := S.Records[S.Keys[j]]
	return a.Score() > b.Score()
}

// Swap swaps the elements with indexes i and j.
func (S RecordSorter) Swap(i, j int) {
	S.Keys[i], S.Keys[j] = S.Keys[j], S.Keys[i]
}

func (R *Result) SortKeys() {
	// Populate the list of Keys:
	for player, record := range R.Records {
		for opponent, _ := range record {
			R.OrderedKeys[player] = append(R.OrderedKeys[player], opponent)
		}
		// Sort the list of Keys based on score:
		data := RecordSorter{
			Keys:    R.OrderedKeys[player],
			Records: R.Records[player],
		}
		sort.Sort(data)
	}
}
