package models

// NewGameList generates a list of games that will need to be played in the tournament.
// Order matters.
func NewGameList(tournamentID Id, s Settings) []*Game {
	// TODO: Carousel
	l := make([]*Game, 0)
	players := append(s.Contestants, s.Opponents...)
	roundNum := 1
	for seat := 0; seat < len(s.Contestants); seat++ {
		for opponent := seat + 1; opponent < len(players); opponent++ {
			for round := 0; round < s.Rounds; round++ {
				m1 := NewGame(tournamentID, s.TimeControl, s.Contestants[seat], players[opponent])
				m1.Round = roundNum
				roundNum++
				m2 := NewGame(tournamentID, s.TimeControl, players[opponent], s.Contestants[seat])
				m2.Round = roundNum
				roundNum++
				l = append(l, []*Game{m1, m2}...)
			}
		}
	}
	setGameOpenings(l, s)
	return l
}
