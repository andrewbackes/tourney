package models

type Engine struct {
	Id         Id     `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Os         string `json:"os"`
	URL        string `json:"url,omitempty"`
	Executable string `json:"executable,omitempty"`
	FilePath   string `json:"filepath,omitempty"`
}
