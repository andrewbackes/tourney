// Package master is a service that facilitates communication between workers and servers.
package master

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/andrewbackes/tourney/data/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrBadStatus     error = errors.New("recieved non-ok status code from server")
	ErrNoGames       error = errors.New("no games found")
	ErrNoTournaments error = errors.New("no tournaments found")
)

type MasterService struct {
	url string
}

func New(url string) *MasterService {
	return &MasterService{url: url}
}

func (m *MasterService) NextGame() (*models.Game, error) {
	tid, err := m.nextTournament()
	if err != nil {
		return &models.Game{}, err
	}
	log.Debug("Tournament: ", tid)
	gid, err := m.nextGame(tid)
	if err != nil {
		return &models.Game{}, err
	}
	return m.GetGame(tid, gid)
}

func (m *MasterService) UpdateGame(g *models.Game) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(g)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, m.url+"/tournaments/"+g.TournamentId.String()+"/games/"+g.Id.String(), b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		log.Debug(r.StatusCode)
		return ErrBadStatus
	}

	return nil
}

func (m *MasterService) UpdatePosition(tid, gid models.Id, p *models.Position) error {
	return nil
}

func (m *MasterService) nextTournament() (models.Id, error) {
	paths := []string{"/tournaments?status=running", "/tournaments?status=pending"}
	for _, path := range paths {
		var t []models.Tournament
		err := m.getJSON(path, &t)
		if err != nil {
			return "", err
		}
		if len(t) > 0 {
			return t[0].Id, nil
		}
	}
	return "", ErrNoTournaments
}

func (m *MasterService) nextGame(tid models.Id) (models.Id, error) {
	path := "/tournaments/" + tid.String() + "/games?status=pending"
	var g []models.Game
	err := m.getJSON(path, &g)
	if err != nil {
		return "", err
	}
	if len(g) > 0 {
		return g[0].Id, nil
	}
	return "", ErrNoGames
}

func (m *MasterService) GetGame(tid, gid models.Id) (*models.Game, error) {
	r, err := http.Get(m.url + "/tournaments/" + tid.String() + "/games/" + gid.String())
	log.Debug(r)
	if err != nil {
		return &models.Game{}, err
	}
	defer r.Body.Close()
	g := &models.Game{}
	err = json.NewDecoder(r.Body).Decode(g)
	return g, err
}

func (m *MasterService) getJSON(path string, obj interface{}) error {
	r, err := http.Get(m.url + path)
	log.Debug(r)
	if err == nil {
		defer r.Body.Close()
		if r.StatusCode != http.StatusOK {
			return ErrBadStatus
		}
		err = json.NewDecoder(r.Body).Decode(&obj)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
