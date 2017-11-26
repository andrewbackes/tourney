package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (m *MemDB) ReadEngines() []*models.Engine {
	result := make([]*models.Engine, 0)
	m.engines.Range(func(id, engine interface{}) bool {
		result = append(result, engine.(*models.Engine))
		return true
	})
	return result
}

func (m *MemDB) CreateEngine(e *models.Engine) {
	_, exists := m.engines.Load(e.Id())
	if !exists {
		m.engines.Store(e.Id(), e)
	}
}

func (m *MemDB) restoreEngines() {
	/*
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
	*/
}
