package main

import (
	"testing"
)

func mockBook() *Book {
	p := map[string]BookPosition{
		"a": BookPosition{Weight: 1},
		"b": BookPosition{Weight: 3},
		"c": BookPosition{Weight: 2},
		"d": BookPosition{Weight: 5},
		"e": BookPosition{Weight: 4},
	}
	i := []string{"a", "b", "c", "d", "e"}
	return &Book{
		Positions: []map[string]BookPosition{p},
		Iterator:  [][]string{i},
	}

}

func TestSortByOccurrence(t *testing.T) {
	b := mockBook()
	b.sortByOccurrence()
	expected := []string{"a", "c", "b", "e", "d"}
	for i, _ := range b.Iterator[0] {
		if b.Iterator[0][i] != expected[i] {
			t.Error("expected:", expected, "but got", b.Iterator[0])
		}
	}
}

func TestRandomize(t *testing.T) {
	b := mockBook()
	b.Randomize(1234567890)
	//expected := []int{3, 1, 1, 0, 0}
	t.Log(b.Iterator[0])
}
