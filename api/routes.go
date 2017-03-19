package api

type Method int

const (
	Post Method = iota
	Get
	Put
	Patch
	Delete
	Options
)

type Endpoint struct {
	Method Method
	Path   string
}

type CommandLine struct {
}

type Route struct {
	Endpoint Endpoint
}

/*

Endpoints:

	/tournaments
	/tournaments/<tourney-id>
	/tournaments/<tourney-id>/games
	/tournaments/<tourney-id>/games/<game-id>
	/tournaments/<tourney-id>/games/<game-id>/plies/<#>
	/tournaments/<tourney-id>/engines

*/
