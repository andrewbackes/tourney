package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Id string

func NewId() Id {
	return Id(bson.NewObjectId().Hex())
}

func (i Id) String() string {
	return string(i)
}
