package models

import (
	"time"
)

type Tournament struct {
	Id           Id                `json:"id"`
	Name         string            `json:"name"`
	CreationDate time.Time         `json:creationDate`
	Tags         map[string]string `json:"tags"`
	Status       Status            `json:"status"`
	Settings     Settings          `json:"settings"`
	Games        []*Game           `json:"games"`
}

type CollapsedTournament struct {
	Id           Id                `json:"id"`
	Name         string            `json:"name"`
	CreationDate time.Time         `json:creationDate`
	Tags         map[string]string `json:"tags"`
	Status       Status            `json:"status"`
	Settings     Settings          `json:"settings"`
}

func CollapseTournament(t *Tournament) *CollapsedTournament {
	return &CollapsedTournament{
		Id:           t.Id,
		Name:         t.Name,
		CreationDate: t.CreationDate,
		Tags:         t.Tags,
		Status:       t.Status,
		Settings:     t.Settings,
	}
}
