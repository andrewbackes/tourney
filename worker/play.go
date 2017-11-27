package worker

import (
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
	log.Info("Playing ", m.TournamentId, "/", m.Id, " - Round ", m.Round)
	engs, err := startEngines(m.Contestants)
	if err != nil {
		panic(err)
	}
	defer closeEngines(engs)
	g := newGame(m)
	engineOutput := make(chan []byte, channelBufferSize)
	status := game.InProgress
	for color := piece.White; status == game.InProgress; color = piece.Color((color + 1) % 2) {
		e := engs[color]
		start := time.Now()
		info, err := e.BestMove(g, engineOutput)
		dur := time.Now().Sub(start)

		w.client.UpdateGame(m)
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
				m.Positions = append(m.Positions, modelPosition(g.Position()))
				m.Positions[len(m.Positions)-1].Analysis = toArray(engineOutput)
				w.client.UpdateGame(m)
			} else {
				status = map[piece.Color]game.GameStatus{piece.White: game.WhiteResigned, piece.Black: game.BlackResigned}[color]
			}
		}
	}
	log.Info(status)
	m.Status = models.Complete
	m.Result = result(status)
	m.EndingCondition = endingCondition(status)
	w.client.UpdateGameWithRetry(m)
	log.Info("Completed game ", m.TournamentId, "/", m.Id, " - Round ", m.Round)
}

func toArray(engineOutput chan []byte) []string {
	analysis := []string{}
	for {
		select {
		case output := <-engineOutput:
			log.Info(string(output))
			analysis = append(analysis, string(output))
		default:
			if len(engineOutput) == 0 {
				return analysis
			}
		}
	}
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
	return map[game.GameStatus]models.EndingCondition{
		game.WhiteCheckmated:      models.Checkmate,
		game.BlackCheckmated:      models.Checkmate,
		game.BlackIllegalMove:     models.IllegalMove,
		game.WhiteIllegalMove:     models.IllegalMove,
		game.BlackResigned:        models.Resignation,
		game.WhiteResigned:        models.Resignation,
		game.BlackTimedOut:        models.OutOfTime,
		game.WhiteTimedOut:        models.OutOfTime,
		game.Stalemate:            models.Stalemate,
		game.InsufficientMaterial: models.InsufficientMaterial,
		game.FiftyMoveRule:        models.FiftyMoveRule,
		game.Threefold:            models.Threefold,
	}[status]
}

func startEngines(e map[piece.Color]models.Engine) (map[piece.Color]*engines.UCIEngine, error) {
	if e[piece.White].ExecPath() == "" || e[piece.Black].ExecPath() == "" {
		panic("engine file path not set")
	}
	w, err := engines.NewUCIEngine(e[piece.White].ExecPath())
	if err != nil {
		return nil, err
	}
	b, err := engines.NewUCIEngine(e[piece.Black].ExecPath())
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
	err := w.client.UpdateGame(g)
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
