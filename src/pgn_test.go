package main

import (
	"testing"
)

func TestSplitTag(t *testing.T) {
	tests := [][]string{
		// test string, expected key, expected value
		{"[Event \"This is a test\"]", "Event", "This is a test"},
		{"[WhiteElo \"2750\"]", "WhiteElo", "2750"},
		{"[Site \"??\"]", "Site", "??"},
		{"[Date \"1.1.2015\"]", "Date", "1.1.2015"},
	}
	for _, test := range tests {
		arg := []byte(test[0])
		key := []byte(test[1])
		val := []byte(test[2])
		if k, v := SplitTag(arg); string(k) != string(key) || string(v) != string(val) {
			t.Error("(", k, ",", v, ") != (", key, ",", val, ")")
		}
	}
}

func TestRemoveComments(t *testing.T) {
	// example: 1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3 ;asdasdasdasd"
	test := []string{
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3 ;asdasdasdasd",
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3",
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd}",
	}
	answer := []string{
		"1. e2e4 d7d5 2. b1c3 f7f5 3. a2a3",
		"1. e2e4 d7d5 2. b1c3 f7f5 3. a2a3",
		"1. e2e4 d7d5 2. b1c3 f7f5",
	}
	for i, _ := range test {
		if result := RemoveComments([]byte(test[i])); string(result) != answer[i] {
			t.Error("'", string(result), "' != '", answer[i], "'")
		}
	}
}

/*
func TestRemoveNumbering(t *testing.T) {
	test := []string{
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3 ;asdasdasdasd",
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3",
		"1. e2e4 d7d5 2. b1c3 f7f5 {asd asd}",
		"1.e4 c5 2.Nf3 d6 3.d4 cxd4 4.Nxd4 Nf6",
	}
	answer := []string{
		"e2e4 d7d5 b1c3 f7f5 {asd asd} a2a3 ;asdasdasdasd",
		"e2e4 d7d5 b1c3 f7f5 {asd asd} a2a3",
		"e2e4 d7d5 b1c3 f7f5 {asd asd}",
		"e4 c5 Nf3 d6 d4 cxd4 Nxd4 Nf6",
	}
	for i, _ := range test {
		if result := RemoveNumbering([]byte(test[i])); string(result) != answer[i] {
			t.Error("'", string(result), "' != '", answer[i], "'")
		}
	}
}
*/

/*
func TestparseLineInPGN(t *testing.T) {
	//pgntext := "[Event \"?\"]\n[Site \"?\"]\n[Date \"2010.06.21\"]\n[Round \"?\"]\n[White \"Alekhine's Defense\"]\n[Black \"?\"]\n[Result \"0-1\"]\n[ECO \"B03\"]\n[PlyCount \"7\"]\n[EventDate \"2010.??.??\"]\n\n1. e4 Nf6 2. e5 Nd5 3. d4 d6 4. Nf3 0-1"

	// Tests if games are getting marked as 'complete' correctly:

	// Situation 1:
	readingmoves := false
	currentGame := NewPGNGame()
	currentGame.Tags["Event"] = "?"
	line := "1. e4 Nf6 2. e5 Nd5 3. d4 d6 4. Nf3 0-1"
	completed, _ := parseLineInPGN(line, &readingmoves, &currentGame)
	if completed {
		t.Error("Situation 1: not expecting 'completed'")
	}

	// Situation 2:
	line = "[Event \"Second test\"]"
	completed, _ = parseLineInPGN(line, &readingmoves, &currentGame)
	if !completed {
		t.Error("Situation 2: expected 'completed'")
	}

}
*/

func TestValuesMatch(t *testing.T) {
	tests := [][]string{
		{"2701", ">2700"}, {"2699", ">2700"},
		{"2699", "<2700"}, {"2800", "<2700"},
		{"apples", "apples"}, {"apples", "oranges"},
		{"a", "=a"}, {"a", "=b"},
	}
	answers := []bool{
		true, false,
		true, false,
		true, false,
		true, false,
	}
	for i, test := range tests {

		if v := ValuesMatch(test[0], test[1]); v != answers[i] {
			t.Error("valuesMatch(", test[0], ",", test[1], ") != ", answers[i])
		}
	}
}

/*
	// **********************************
	filters := []PGNFilter{
		//{Tag: "FEN", Value: "rn1qk2r/pb2ppbp/1p1p1np1/8/2PQ4/2N2NP1/PP2PPBP/R1B2RK1 w kq - 0 9"},
		{Tag: "WhiteElo", Value: ">2699"},
		{Tag: "BlackElo", Value: ">2699"},
		{Tag: "Result", Value: "1/2-1/2"},
	}

	list, _ := ReadPGN("/Users/Andrew/Documents/millionbase-2.22.pgn", filters)
	//list, _ := ReadPGN("default.pgn", filters)
	//list, _ := ReadPGN("Alekhine.pgn", filters)
	for _, g := range *list {
		fmt.Println(g.Tags["White"], g.Tags["WhiteElo"], "vs", g.Tags["Black"], g.Tags["BlackElo"])
	}
	fmt.Println("Count:", len(*list))
	return
*/
// **********************************
