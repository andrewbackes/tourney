package models

type Tournament struct {
	Id       Id                `json:"id"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
	Games    []*Game           `json:"games"`
}

type CollapsedTournament struct {
	Id       Id                `json:"id"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
}

func CollapseTournament(t *Tournament) *CollapsedTournament {
	return &CollapsedTournament{
		Id:       t.Id,
		Tags:     t.Tags,
		Status:   t.Status,
		Settings: t.Settings,
	}
}
