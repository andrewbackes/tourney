package util

import (
	"archive/zip"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src string) (string, error) {
	unzippedDir := filepath.Dir(src)

	/*
		if _, err := os.Stat(unzippedDir); !os.IsNotExist(err) {
			log.Info(unzippedDir, " already exists. Skipping unzip.")
			return unzippedDir, nil
		}
	*/

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
		fpath := filepath.Join(unzippedDir, f.Name)
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
