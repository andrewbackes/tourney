package api

import (
	"testing"
	"bytes"
)

func TestWriteJSON(t *testing.T) {
	var b bytes.Buffer
	m  := make(map[string]string)
	m["mykey"] = "myval"
	writeJSON(m, &b)
	expected := `{
  "mykey": "myval"
}
`
	actual := b.String() 
	if actual !=  expected {
		t.Error("Got:", actual, "\nWanted:", expected)
	}
}