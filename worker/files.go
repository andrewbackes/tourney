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

const storageLocation = "tourney_storage/downloadedEngineFiles"

func (w *Worker) getEngines(g *models.Game) {
	for color, engine := range g.Contestants {
		dl, err := download(engine)
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
			err = os.Chmod(engine.FilePath, 755)
			if err != nil {
				panic("Could not set permissions on " + engine.FilePath + " executable")
			}
			log.Info("Found executable ", engine.FilePath)
		} else if os.IsNotExist(err) {
			panic(engine.FilePath + " does not exist")
		}
		g.Contestants[color] = engine
	}
}

func download(engine models.Engine) (string, error) {
	url := urlOf(engine)
	absPath := filepath.Join(storageLocation, filenameOf(engine))
	err := os.MkdirAll(storageLocation, os.ModePerm)
	if err != nil {
		panic(err)
	}
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
	if err != nil || response.StatusCode > 399 {
		log.Error("Error while downloading ", url, " - ", response.StatusCode, " - ", err)
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

func urlOf(engine models.Engine) string {
	if engine.URL != "" {
		return engine.URL
	}
	return getAPIURL() + "/engineFiles/" + engine.Name + "/" + engine.Version + "/" + engine.Os
}

func filenameOf(engine models.Engine) string {
	if engine.URL != "" {
		terms := strings.Split(engine.URL, "/")
		return terms[len(terms)-1]
	}
	return engine.Id()
}

func getAPIURL() string {
	if os.Getenv("API_URL") != "" {
		return os.Getenv("API_URL")
	}
	return "http://api.tourney.aback.es:9090/api/v2"
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
