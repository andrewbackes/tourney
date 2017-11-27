package worker

import (
	"errors"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (w *Worker) getEngines(g *models.Game) {
	for color, engine := range g.Contestants {
		if _, err := os.Stat(engine.ExecPath()); !os.IsNotExist(err) {
			log.Info(engine.ExecPath(), " already exists. Skipping download.")
			return
		}
		absPath, err := download(engine)
		if err != nil {
			panic(err)
		}
		if strings.HasSuffix(absPath, ".zip") {
			_, err := util.Unzip(absPath)
			if err != nil {
				panic(err)
			}
		} else if strings.HasSuffix(absPath, ".tar.gz") || strings.HasSuffix(absPath, ".tgz") {
			panic(errors.New("tar files not yet supported"))
		} else if strings.HasSuffix(absPath, ".rar") {
			panic(errors.New("rar files not yet supported"))
		}
		if _, err := os.Stat(engine.ExecPath()); err == nil {
			err = os.Chmod(engine.ExecPath(), 755)
			if err != nil {
				panic("Could not set permissions on " + engine.ExecPath() + " executable")
			}
			log.Info("Found executable ", engine.ExecPath())
		} else if os.IsNotExist(err) {
			panic(engine.ExecPath() + " does not exist")
		}
		g.Contestants[color] = engine
	}
}

func download(engine models.Engine) (string, error) {
	terms := strings.Split(engine.URL, "/")
	filename := terms[len(terms)-1]
	absPath := filepath.Join(engine.DirPath(), filename)
	err := os.MkdirAll(engine.DirPath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.Info("Downloading ", engine.URL, " to ", absPath)

	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		log.Info(absPath, " already exists. Skipping download.")
		return absPath, nil
	}

	output, err := os.Create(absPath)
	if err != nil {
		log.Error("Error while creating ", absPath, " - ", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(engine.URL)
	if err != nil || response.StatusCode > 399 {
		log.Error("Error while downloading ", engine.URL, " - ", response.StatusCode, " - ", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Error("Error while downloading ", engine.URL, " - ", err)
		return "", err
	}

	log.Info(n, " bytes downloaded.")
	return absPath, nil
}
