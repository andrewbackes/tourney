package memdb

import (
	"encoding/json"
	"github.com/andrewbackes/tourney/data/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		m.persistEngine(e)
	}
}

func (m *MemDB) persistEngine(e *models.Engine) {
	if !m.persisted() {
		return
	}
	engineDir := filepath.Join(m.backupDir, "engines")
	err := os.MkdirAll(engineDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	engineJSON := filepath.Join(engineDir, e.Id()+".json")
	f, err := os.Create(engineJSON)
	if err != nil {
		panic(err)
	}
	log.Info("Persisting engine", engineJSON)
	err = json.NewEncoder(f).Encode(e)
	if err != nil {
		panic(err)
	}
}

func (m *MemDB) restoreEngines() {
	root := filepath.Join(m.backupDir, "engines")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Info("No engines found in ", root)
		return
	}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			log.Debug("Restoring engine ", file.Name())
			f, err := os.Open(filepath.Join(root, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			e := &models.Engine{}
			err = json.NewDecoder(f).Decode(e)
			if err != nil {
				log.Fatal(err)
			}
			m.engines.Store(e.Id(), e)
			f.Close()
		}
	}
}
