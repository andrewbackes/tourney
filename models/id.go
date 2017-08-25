package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Id string

func NewId() Id {
	return Id(bson.NewObjectId())
}
