package models

type Tournament struct {
	Id       Id                `json:"id"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
	Summary  Summary           `json:"summary"`
	Games    []*Game           `json:"games"`
}
