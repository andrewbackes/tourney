package structures

import (
	"testing"
)

func TestPlaceholder(t *testing.T) {
}

/*
func TestCarousel(t *testing.T) {
	trn := NewTournament()
	trn.Carousel = true
	trn.Rounds = 3
	trn.Contestants = []Engine{
		Engine{Id: "0"},
		Engine{Id: "1"},
		Engine{Id: "2"},
	}
	gs := trn.GenerateGames()
	expected := [][]string{
		{"30", "31"},
		{"30", "32"},
		{"31", "30"},
		{"32", "30"},
		{"30", "31"},
		{"30", "32"},
	}

	for i, g := range gs {
		if g.Tags["WhiteId"] != expected[i][0] {
			t.Error(g.Tags["WhiteId"], "not", expected[i][0])
		}
		if g.Tags["BlackId"] != expected[i][1] {
			t.Error(g.Tags["BlackId"], "not", expected[i][1])
		}

	}
}

func TestNonCarousel(t *testing.T) {
	trn := NewTournament()
	trn.Carousel = false
	trn.Rounds = 3
	trn.Contestants = []Engine{
		Engine{Id: "0"},
		Engine{Id: "1"},
		Engine{Id: "2"},
	}
	gs := trn.GenerateGames()
	expected := [][]string{
		{"30", "31"},
		{"31", "30"},
		{"30", "31"},
		{"30", "32"},
		{"32", "30"},
		{"30", "32"},
	}

	for i, g := range gs {
		t.Log("Got:", g.Tags["WhiteId"], "vs", g.Tags["BlackId"])
		if g.Tags["WhiteId"] != expected[i][0] {
			t.Error(g.Tags["WhiteId"], "not", expected[i][0])
		}
		if g.Tags["BlackId"] != expected[i][1] {
			t.Error(g.Tags["BlackId"], "not", expected[i][1])
		}

	}
}

func TestCarouselMultipleSeats(t *testing.T) {
	trn := NewTournament()
	trn.Carousel = true
	trn.Rounds = 3
	trn.TestSeats = 3
	trn.Contestants = []Engine{
		Engine{Id: "0"},
		Engine{Id: "1"},
		Engine{Id: "2"},
	}
	gs := trn.GenerateGames()
	expected := [][]string{
		{"30", "31"},
		{"30", "32"},
		{"31", "30"},
		{"32", "30"},
		{"30", "31"},
		{"30", "32"},
		{"31", "32"},
		{"32", "31"},
		{"31", "32"},
	}

	t.Log(expected)
	for i, g := range gs {
		t.Log("Got:", g.Tags["WhiteId"], "vs", g.Tags["BlackId"])

		if g.Tags["WhiteId"] != expected[i][0] {
			t.Error(g.Tags["WhiteId"], "not", expected[i][0])
		}
		if g.Tags["BlackId"] != expected[i][1] {
			t.Error(g.Tags["BlackId"], "not", expected[i][1])
		}

	}

}
*/
