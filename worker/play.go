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

const (
	channelBufferSize = 256
)

func (w *Worker) play(m *models.Game) {
	fmt.Println(m)
	engs, err := startEngines(m.Contestants)
	if err != nil {
		panic(err)
	}
	defer closeEngines(engs)
	g := newGame(m)
	engineOutput := make(chan []byte, channelBufferSize)
	positionFeed := make(chan *position.Position, channelBufferSize)
	done := make(chan struct{})
	go w.positionUpdater(m, engineOutput, positionFeed, done)
	status := game.InProgress
	for color := piece.White; status == game.InProgress; color = piece.Color((color + 1) % 2) {
		e := engs[color]
		start := time.Now()
		info, err := e.BestMove(g, engineOutput)
		dur := time.Now().Sub(start)
		if err == engines.ErrTimedOut {
			status = map[piece.Color]game.GameStatus{piece.White: game.WhiteTimedOut, piece.Black: game.BlackTimedOut}[color]
		} else {
			bm := move.Parse(info.BestMove)
			bm.Duration = dur
			log.Debug(info)
			log.Info(bm)
			if bm != move.Null {
				status, err = g.MakeMove(bm)
				if err != nil {
					panic(err)
				}
				positionFeed <- g.Position()
			} else {
				status = map[piece.Color]game.GameStatus{piece.White: game.WhiteResigned, piece.Black: game.BlackResigned}[color]
			}
		}
	}
	log.Info(status)
	m.Status = models.Complete
	m.Result = result(status)
	m.EndingCondition = endingCondition(status)
	w.master.UpdateGame(m)
}

func result(status game.GameStatus) models.Result {
	if status&game.WhiteWon != 0 {
		return models.White
	} else if status&game.BlackWon != 0 {
		return models.Black
	}
	return models.Draw
}

func endingCondition(status game.GameStatus) models.EndingCondition {
	switch status {
	case game.WhiteCheckmated:
		return models.Checkmate
	case game.BlackCheckmated:
		return models.Checkmate
	case game.BlackIllegalMove:
		return models.IllegalMove
	case game.WhiteIllegalMove:
		return models.IllegalMove
	case game.BlackResigned:
		return models.Resignation
	case game.WhiteResigned:
		return models.Resignation
	case game.BlackTimedOut:
		return models.OutOfTime
	case game.WhiteTimedOut:
		return models.OutOfTime
	case game.Stalemate:
		return models.Stalemate
	case game.InsufficientMaterial:
		return models.InsufficientMaterial
	case game.FiftyMoveRule:
		return models.FiftyMoveRule
	case game.Threefold:
		return models.Threefold
	default:
		return models.Error
	}
}

func (w *Worker) positionUpdater(m *models.Game, engineOutput chan []byte, positionFeed chan *position.Position, done chan struct{}) {
	for {
		select {
		case output := <-engineOutput:
			log.Info(string(output))
			m.Positions[len(m.Positions)-1].Analysis = append(m.Positions[len(m.Positions)-1].Analysis, string(output))
			w.master.UpdateGame(m)
		case pos := <-positionFeed:
			//m.Positions[len(m.Positions)-1].LastAnalysis = info.Analysis
			m.Positions = append(m.Positions, modelPosition(pos))
			w.master.UpdateGame(m)
		case <-done:
			if len(engineOutput) == 0 && len(positionFeed) == 0 {
				return
			}
		}
	}
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

func modelPosition(p *position.Position) models.Position {
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
		LastMove: p.LastMove,
	}
}
