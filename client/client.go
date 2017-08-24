package client

import (
	"fmt"
	"github.com/andrewbackes/tourney/helpers"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

var (
	NoObjectId = bson.ObjectId("")
)

type Client struct {
	apiBaseUrl string
	stop       chan struct{}
}

func New(hostname string) *Client {
	return &Client{
		apiBaseUrl: "http://" + hostname + "/api/v2",
		stop:       make(chan struct{}),
	}
}

func (c *Client) getTournament() bson.ObjectId {
	resp, err := http.Get(fmt.Sprintf("%s/tournaments?completed=false", c.apiBaseUrl))
	if err != nil {
		panic(err)
	}
	var ts []structures.Tournament
	helpers.ReadJSON(resp.Body, &ts)
	resp.Body.Close()
	if ts == nil || len(ts) == 0 {
		return NoObjectId
	}
	return ts[0].Id
}

func (c *Client) getNextGame(tournamentId bson.ObjectId) *structures.Game {
	resp, err := http.Get(fmt.Sprintf("%s/tournaments/%s/gameQueue/next", c.apiBaseUrl, tournamentId.Hex()))
	if err != nil {
		panic(err)
	}
	var g structures.Game
	helpers.ReadJSON(resp.Body, &g)
	if g.Id == "" {
		fmt.Println("game id is blank")
		panic(&g)
	}
	return &g
}

func (c *Client) Stop() {
	if c.stopped() {
		fmt.Println("Client already stopped.")
	} else {
		c.Stop()
	}
}

func (c *Client) ConnectAndPlay() {
	for !c.stopped() {
		tid := c.getTournament()
		if tid != NoObjectId {
			fmt.Printf("Tournament set to %s\n", tid.Hex())
			g := c.getNextGame(tid)
			c.play(g)
		} else {
			fmt.Println("No tournaments to play.")
			time.Sleep(10 * time.Second)
		}
	}
}

func (c *Client) stopped() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}
