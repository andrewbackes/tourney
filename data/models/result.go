package models

type Result int

const (
	Incomplete Result = iota
	White
	Black
	Draw
)
