package worker

import (
	"archive/zip"
	"errors"
	"github.com/andrewbackes/tourney/data/models"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const storageLocation = "/Users/Andrew/tourney_engines/"

func (w *Worker) getEngines(g *models.Game) {
	for color, engine := range g.Contestants {
		dl, err := downloadFromURL(engine.URL)
		if err != nil {
			panic(err)
		}
		engine.FilePath = dl
		if strings.HasSuffix(dl, ".zip") {
			path, err := unzip(dl)
			if err != nil {
				panic(err)
			}
			engine.FilePath = filepath.Join(path, engine.Executable)

		} else if strings.HasSuffix(dl, ".tar.gz") || strings.HasSuffix(dl, ".tgz") {
			path, err := untar(dl)
			if err != nil {
				panic(err)
			}
			engine.FilePath = filepath.Join(path, engine.Executable)
		} else if strings.HasSuffix(dl, ".rar") {
			path, err := unrar(dl)
			if err != nil {
				panic(err)
			}
			engine.FilePath = filepath.Join(path, engine.Executable)
		}
		if _, err := os.Stat(engine.FilePath); err == nil {
			log.Info("Found executable ", engine.FilePath)
		} else if os.IsNotExist(err) {
			panic(engine.FilePath + " does not exist")
		}
		g.Contestants[color] = engine
	}
}

func downloadFromURL(url string) (string, error) {
	tokens := strings.Split(url, "/")
	absPath := filepath.Join(storageLocation, tokens[len(tokens)-1])
	log.Info("Downloading ", url, " to ", absPath)

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

	response, err := http.Get(url)
	if err != nil {
		log.Error("Error while downloading ", url, " - ", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Error("Error while downloading ", url, " - ", err)
		return "", err
	}

	log.Info(n, " bytes downloaded.")
	return absPath, nil
}

func unzip(src string) (string, error) {
	unzippedDir := src[:len(src)-len(".zip")]
	if _, err := os.Stat(unzippedDir); !os.IsNotExist(err) {
		log.Info(unzippedDir, " already exists. Skipping unzip.")
		return unzippedDir, nil
	}

	log.Info("Unzipping ", src, " to ", unzippedDir)
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return unzippedDir, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return unzippedDir, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(storageLocation, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return unzippedDir, err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return unzippedDir, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return unzippedDir, err
			}

		}
	}
	log.Info("Unzipped: ", filenames)
	return unzippedDir, nil
}

func unrar(src string) (string, error) {
	return src, errors.New("rar files not yet supported")
}

func untar(src string) (string, error) {
	return src, errors.New("tar files not yet supported")
}
