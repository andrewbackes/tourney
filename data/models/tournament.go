package models

type Tournament struct {
	Id       Id                `json:"id"`
	Tags     map[string]string `json:"tags"`
	Status   Status            `json:"status"`
	Settings Settings          `json:"settings"`
	Summary  Summary           `json:"summary"`
	Games    []*Game           `json:"games"`
}

// NewGameList generates a list of games that will need to be played in the tournament.
// Order matters.
func NewGameList(tournamentID Id, s Settings) []*Game {
	// TODO: Carousel
	l := make([]*Game, 0)
	for seat := 0; seat < s.TestSeats; seat++ {
		for opponent := seat + 1; opponent < len(s.Engines); opponent++ {
			for round := 0; round < s.Rounds; round++ {
				m1 := NewGame(tournamentID, s.TimeControl, s.Engines[seat], s.Engines[opponent])
				m2 := NewGame(tournamentID, s.TimeControl, s.Engines[opponent], s.Engines[seat])
				l = append(l, []*Game{m1, m2}...)
			}
		}
	}
	/*
		for i := 0; i < s.TestSeats; i++ {
			if s.Carousel {
				for round := 0; round < s.Rounds; round = round + []int{2, 1}[s.Rounds%2] {
					for engine := i + 1; engine < len(s.Engines); engine++ {
						g1 := &Game{
							TournamentId: tournamentID,
							Id:           NewId(),
							TimeControl:  s.TimeControl,
						}
						g1.Contestants[piece.Color(round%2)] = s.Engines[i]
						g1.Contestants[piece.Color((round+1)%2)] = s.Engines[engine]
						l = append(l, g1)
						if s.Rounds%2 == 0 {
							g2 := &Game{
								TournamentId: tournamentID,
								Id:           NewId(),
								TimeControl:  s.TimeControl,
							}
							g2.Contestants[piece.Color(round%2)] = s.Engines[engine]
							g2.Contestants[piece.Color((round+1)%2)] = s.Engines[i]
							l = append(l, g2)
						}
					}
				}
			} else {
				for engine := i + 1; engine < len(s.Engines); engine++ {
					for round := 0; round < s.Rounds; round++ {
						g1 := &Game{
							TournamentId: tournamentID,
							Id:           NewId(),
							TimeControl:  s.TimeControl,
						}
						g1.Contestants[piece.Color(round%2)] = s.Engines[i]
						g1.Contestants[piece.Color((round+1)%2)] = s.Engines[engine]
						l = append(l, g1)
					}
				}
			}
		}
	*/
	return l
}
