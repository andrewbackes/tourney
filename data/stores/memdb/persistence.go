package memdb

import (
	"encoding/json"
	"github.com/andrewbackes/tourney/data/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func (m *MemDB) restore() {
	log.Info("Loading persistent data from ", m.backupDir)
	m.restoreTournaments()
	m.restoreEngines()
}

func (m *MemDB) restoreTournaments() {
	root := filepath.Join(m.backupDir, "tournaments")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Info("No tournaments found in ", root)
		return
	}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			log.Debug("Restoring tournament ", file.Name())
			f, err := os.Open(filepath.Join(root, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			t := &models.Tournament{}
			err = json.NewDecoder(f).Decode(t)
			if err != nil {
				log.Fatal(err)
			}
			m.tournaments.Store(t.Id, t)
			m.locks.Store(t.Id, &sync.Mutex{})
			f.Close()
			gameDir := filepath.Join(root, strings.TrimSuffix(file.Name(), ".json"))
			gameFiles, err := ioutil.ReadDir(gameDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, gameFile := range gameFiles {
				gf, err := os.Open(filepath.Join(gameDir, gameFile.Name()))
				if err != nil {
					log.Fatal(err)
				}
				g := models.Game{}
				err = json.NewDecoder(gf).Decode(&g)
				if err != nil {
					log.Fatal(err)
				}
				m.games.Store(g.Id, &g)
				m.locks.Store(g.Id, &sync.Mutex{})
				t, ok := m.tournaments.Load(g.TournamentId)
				if !ok {
					log.Fatal("can not add games to a tournament that does not exist")
				}
				t.(*models.Tournament).Games = append(t.(*models.Tournament).Games, &g)
				gf.Close()
			}
		}
	}
}
