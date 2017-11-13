package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
)

func (s *Service) ReadGame(tid, gid models.Id) (*models.Game, error) {
	g, err := s.store.ReadGame(tid, gid)
	log.Debug("Read Game ", *g)
	if err == stores.ErrNotFound {
		return nil, ErrNotFound
	}
	return g, nil
}

func (s *Service) ReadGames(tid models.Id, filter func(*models.Game) bool) []*models.Game {
	return s.store.ReadGames(tid, filter)
}

func (s *Service) UpdateGame(g *models.Game) error {
	s.store.UpdateGame(g)
	t, err := s.store.ReadTournament(g.TournamentId)
	if err != nil {
		log.Error(err)
		return err
	}
	if g.Status == models.Running {
		if t.Status != models.Running {
			t.Status = models.Running
			s.store.UpdateTournament(t)
		}
	} else if g.Status == models.Complete {
		complete := true
		for _, tg := range t.Games {
			if tg.Status != models.Complete {
				complete = false
				break
			}
		}
		if complete {
			t.Status = models.Complete
		}
		t.Summary = models.NewSummary(t.Settings.Contestants, t.Games)
		s.store.UpdateTournament(t)
	}
	return nil
}

/*
func (s *Service) AddPosition(tid, gid models.Id, p models.Position) error {
	a := strings.Split(p.FEN, " ")
	if len(a) < 6 {
		return ErrMalformed
	}
	s.store.AddPosition(tid, gid, p)
	return nil
}

func (s *Service) CompleteGame(tid, gid models.Id) {
	s.store.UpdateStatus(tid, gid, models.Complete)
}

func (s *Service) AssignGame(tid, gid models.Id) {
	s.store.UpdateStatus(tid, gid, models.Running)
}
*/
