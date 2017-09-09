package worker

import (
	"fmt"
	"github.com/andrewbackes/chess/engines"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/chess/position/move"
	"github.com/andrewbackes/tourney/data/models"
	log "github.com/sirupsen/logrus"
	"time"
)

func (w *Worker) play(m *models.Game) {
	fmt.Println(m)
	engines, err := startEngines(m.Contestants)
	if err != nil {
		panic(err)
	}
	defer closeEngines(engines)
	g := newGame(m)
	status := game.InProgress
	for color := piece.White; status == game.InProgress; color = piece.Color((color + 1) % 2) {
		e := engines[color]
		start := time.Now()
		info, err := e.BestMove(g)
		dur := time.Now().Sub(start)
		if err != nil {
			panic(err)
		}
		bm := move.Parse(info.BestMove)
		bm.Duration = dur
		log.Debug(info)
		log.Info(bm)
		status, err = g.MakeMove(bm)
		if err != nil {
			panic(err)
		}
		m.Positions = append(m.Positions, modelPosition(g.Position(), info.Analysis))
		go func() {
			w.master.UpdateGame(m)
		}()
	}
	log.Info(status)
	m.Status = models.Complete
	w.master.UpdateGame(m)
}

func startEngines(e map[piece.Color]models.Engine) (map[piece.Color]*engines.UCIEngine, error) {
	if e[piece.White].FilePath == "" || e[piece.Black].FilePath == "" {
		panic("engine file path not set")
	}
	w, err := engines.NewUCIEngine(e[piece.White].FilePath)
	if err != nil {
		return nil, err
	}
	b, err := engines.NewUCIEngine(e[piece.Black].FilePath)
	if err != nil {
		return nil, err
	}
	w.NewGame()
	b.NewGame()
	return map[piece.Color]*engines.UCIEngine{
		piece.White: w,
		piece.Black: b,
	}, nil
}
func closeEngines(e map[piece.Color]*engines.UCIEngine) {
	w := e[piece.White]
	b := e[piece.Black]
	w.Close()
	b.Close()
}

func newGame(g *models.Game) *game.Game {
	new := game.NewTimedGame(map[piece.Color]game.TimeControl{
		piece.White: g.TimeControl,
		piece.Black: g.TimeControl,
	})
	for _, p := range g.Positions {
		if p.LastMove != move.Null {
			status, err := new.MakeMove(p.LastMove)
			if err != nil {
				panic(err)
			}
			if status != game.InProgress {
				panic("game should not terminate during opening book")
			}
		}
	}
	return new
}

func (w *Worker) claim(g *models.Game) {
	g.Status = models.Running
	err := w.master.UpdateGame(g)
	if err != nil {
		panic(err)
	}
}

func modelPosition(p *position.Position, analysis []map[string]string) models.Position {
	f, err := fen.Encode(p)
	if err != nil {
		log.Error(err)
		return models.Position{}
	}
	return models.Position{
		FEN: f,
		Clocks: map[piece.Color]time.Duration{
			piece.White: p.Clocks[piece.White],
			piece.Black: p.Clocks[piece.Black],
		},
		MovesLeft: map[piece.Color]int{
			piece.White: p.MovesLeft[piece.White],
			piece.Black: p.MovesLeft[piece.Black],
		},
		LastMove:     p.LastMove,
		LastAnalysis: analysis,
	}
}
