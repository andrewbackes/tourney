package models

type Tournament struct {
	Id       Id                `json:"id"`
	Name     string            `json:"name"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
	Games    []*Game           `json:"games"`
}

type CollapsedTournament struct {
	Id       Id                `json:"id"`
	Name     string            `json:"name"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
}

func CollapseTournament(t *Tournament) *CollapsedTournament {
	return &CollapsedTournament{
		Id:       t.Id,
		Name:     t.Name,
		Tags:     t.Tags,
		Status:   t.Status,
		Settings: t.Settings,
	}
}
