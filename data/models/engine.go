package models

import (
	"strings"
)

type Engine struct {
	Name       string `json:"name"`
	Version    string `json:"version,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Os         string `json:"os"`
	URL        string `json:"url,omitempty"`
	Executable string `json:"executable,omitempty"`
	FilePath   string `json:"filepath,omitempty"`
}

func (e *Engine) Id() string {
	return strings.ToLower(e.Name + "-" + e.Version + "-" + e.Os)
}
