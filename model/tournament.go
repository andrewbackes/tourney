package model

import (
	"gopkg.in/mgo.v2/bson"
)

/*
type Status int

type Engine struct {
	FilePath string
}

type Book struct {
	FilePath        string `json:"filePath,omitempty"`
	Depth           int    `json:"depth,omitempty"`
	Randomize       bool   `json:"randomize,omitempty"`
	MirrorPositions bool   `json:"mirrorPositions,omitempty"`
	RepeatPositions bool   `json:"repeatPositions,omitempty"`
}
*/

type Tournament struct {
	Id        bson.ObjectId     `json:"id" bson:"_id"`
	Tags      map[string]string `json:"tags,omitempty" bson:"tags,omitempty"`
	TestSeats int               `json:"testSeats" bson:"testSeats"`
	Carousel  bool              `json:"carousel" bson:"carousel"`
	Rounds    int               `json:"rounds" bson:"rounds"`
}
