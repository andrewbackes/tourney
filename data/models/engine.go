package models

type Engine struct {
	Id       Id
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	URL      string `json:"url,omitempty"`
	FilePath string `json:"filepath,omitempty"`
}
