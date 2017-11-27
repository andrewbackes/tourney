package models

import (
	"github.com/andrewbackes/tourney/util"
	"path/filepath"
	"strings"
)

type Engine struct {
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Os       string `json:"os"`
	// URL to download the engine file(s).
	URL string `json:"url,omitempty"`
	// Executable is the relative path to the engine executable.
	Executable string `json:"executable,omitempty"`
}

func (e *Engine) Id() string {
	return strings.ToLower(e.Name + "-" + e.Version + "-" + e.Os)
}

func (e Engine) ExecPath() string {
	return filepath.Join(util.GetStorageLocation(), "engineFiles", e.Name, e.Version, e.Os, e.Executable)
}

func (e Engine) DirPath() string {
	return filepath.Join(util.GetStorageLocation(), "engineFiles", e.Name, e.Version, e.Os)
}
