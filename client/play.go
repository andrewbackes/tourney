package client

import (
	"bytes"
	"fmt"
	"github.com/andrewbackes/tourney/helpers"
	"github.com/andrewbackes/tourney/model/structures"
	"net/http"
)

func (c *Client) play(g *structures.Game) {
	fmt.Printf("Playing game %s\n", g.Id.Hex())
	fmt.Println(g)
	g.Tags["result"] = "1/2-1/2"
	c.UpdateGameResult(g)
}

func (c *Client) UpdateGameResult(g *structures.Game) {
	fmt.Println("Updating game status")
	var b bytes.Buffer
	helpers.WriteJSON(g, &b)
	url := fmt.Sprintf("%s/tournaments/%s/games/%s", c.apiBaseUrl, g.TournamentId().Hex(), g.Id.Hex())
	req, err := http.NewRequest(http.MethodPatch, url, &b)
	if err != nil {
		panic(err)
	}
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
}
