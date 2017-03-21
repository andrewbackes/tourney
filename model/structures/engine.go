package structures

import (
	"gopkg.in/mgo.v2/bson"
)

type Engine struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Version  string        `json:"version,omitempty" bson:"version,omitempty"`
	Protocol string        `json:"protocol,omitempty" bson:"protocol,omitempty"`
	URL      string        `json:"url,omitempty" bson:"url,omitempty"`
	FilePath string        `json:"filepath,omitempty" bson:"filepath,omitempty"`
}
