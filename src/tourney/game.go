/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

 Module: games
 Description: game object. Lots to explain...

 TODO:
 	-Test the fen save/load by loading a file full of FEN and make sure the
 	same fen comes out after loading.
 	-check state before each modifying function and functions that print to
 	the screen
 	-add support for processing nullmoves.
 	-double check 3 fold and 50 move rules.
 	-optimize 3 fold.

*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	//"runtime"
	"os"
	"strconv"
	"strings"
	"time"
)

type Color uint8

const (
	WHITE   Color = 0
	BLACK   Color = 1
	BOTH    Color = 2
	NEITHER Color = 2
	DRAW    Color = 2
)

const (
	SHORT uint = 0
	LONG  uint = 1
)

type Game struct {

	// Header info (General info, usually needed for pgn):
	Event       string
	Site        string
	Date        string
	Round       int
	Time        int64
	Moves       int64
	Repeating   bool
	Player      [2]Engine // white=0,black=1
	StartingFEN string    // first position out of the book

	// Control info:
	Timer     [2]int64 //the game clock for each side. milliseconds.
	MovesToGo int64    //moves until time control

	// Game run time info:
	Board        Board
	FiftyRule    uint64
	History      []string // FEN of every known position so far.
	EnPassant    uint8
	CastleRights [2][2]bool
	MoveList     []Move //move history
	Completed    bool   // TODO:  change this to reflect how the game ended. time, checkmate, adjunction, etc

	// Post game info:
	Result       Color //WHITE,BLACK,DRAW
	ResultDetail string

	// Logging:
	logFile   *os.File
	logBuffer string
}

/*******************************************************************************

	Primary function. Manages the game itself. Loops until instructed not to.
	Or until the game ends.

*******************************************************************************/

func PlayGame(G *Game) error {
	// Note: opening book is handled in RunTourney()
	G.StartLog()
	fmt.Println("Playing Game...")
	// Start up the engines:
	if err := G.Player[WHITE].Start(&G.logBuffer); err != nil {
		return err
	}
	if err := G.Player[BLACK].Start(&G.logBuffer); err != nil {
		return err
	}

	var state Status = RUNNING
	for state == RUNNING {
		if state == STOPPED {
			// More clean up code here.
			break
		}
		// Play:
		gameover := ExecuteNextTurn(G)
		G.AppendLog()
		if gameover {
			state = STOPPED
			break
		}

		if G.toMove() == WHITE {
			G.MovesToGo -= 1
			if G.MovesToGo == 0 && G.Repeating == true {
				G.resetTimeControl()
			}
		}
	}

	// Stop the engines:
	G.Player[WHITE].Shutdown()
	G.Player[BLACK].Shutdown()

	G.AppendLog()
	G.CloseLog()
	return nil
}

func ExecuteNextTurn(G *Game) bool {
	//Return true to quit game
	color := G.toMove()
	otherColor := []Color{BLACK, WHITE}[color]
	// Tell the engine to set its internal board:
	if err := G.Player[color].Set(G.MoveList); err != nil { //TODO: should probably pass by ref
		G.GameOver(color, err.Error())
		return true
	}
	// Request a move from the engine:
	engineMove, lapsed, err := G.Player[color].Move(G.Timer, G.MovesToGo, color)
	if err != nil {
		G.GameOver(color, err.Error())
		return true
	}
	// Adjust time control:
	G.Timer[color] -= lapsed.Nanoseconds() / 1000000
	if G.Timer[color] < 0 {
		G.GameOver(color, "Out of time. Used "+strconv.FormatInt(lapsed.Nanoseconds()/1000000, 10)+"ms / "+
			strconv.FormatInt(G.Timer[color]+(lapsed.Nanoseconds()/1000000), 10)+"ms.")
		return true
	}
	// Convert the notation from the engines notation to pure coordinate notation
	parsedMove := engineMove
	parsedMove.Algebraic = InternalizeNotation(G, parsedMove.Algebraic)

	// Print the move:
	if color == WHITE {
		fmt.Print(len(G.MoveList)/2+1, ". ")
	}
	fmt.Print(parsedMove.Algebraic, " ")

	// Check legality of move.
	LegalMoves := LegalMoveList(G)
	if !contains(LegalMoves, parsedMove) {
		G.GameOver(color, "Illegal move: "+parsedMove.Algebraic+" (raw: "+engineMove.Algebraic+").")
		return true
	}
	// Adjust the internal board:
	if err = G.MakeMove(parsedMove); err != nil {
		G.GameOver(color, err.Error()) // illegal move
		return true
	}
	// Check:
	check := G.isInCheck(otherColor)
	// Mate:
	oppLegalMoves := LegalMoveList(G)
	if len(oppLegalMoves) == 0 {
		if check {
			//checkmate!
			G.GameOver(otherColor, "Checkmate.")
			return true
		} else {
			//stalemate!
			G.GameOver(NEITHER, "Stalemate.")
			return true
		}
	}
	// 50 Move Draw:
	if FiftyMoveDraw(G) {
		G.GameOver(NEITHER, "50 move rule.")
		return true
	}
	// Insufficient material:
	if InsufficientMaterial(G) {
		G.GameOver(NEITHER, "Insufficient material.")
		return true
	}
	// 3 fold:
	if ThreeFold(G) {
		G.GameOver(NEITHER, "Three fold repetition.")
		return true
	}
	return false
}

func contains(list []Move, move Move) bool {
	for i, _ := range list {
		if move.Algebraic == list[i].Algebraic {
			return true
		}
	}
	return false
}

func FiftyMoveDraw(G *Game) bool {
	return G.FiftyRule == 100
}

func ThreeFold(G *Game) bool {
	// TODO: Can this be optimized? very very slow right now.
	// maybe sort them alphabetically first.
	// should go backwards instead, as many times as in the fiftyMove count.
	for i := 0; i < len(G.History); i++ {
		fen := G.History[i]
		fenSplit := strings.Split(fen, " ")
		fenPrefix := fenSplit[0] + " " + fenSplit[1] + " " + fenSplit[2] + " " + fenSplit[3]
		for j := i + 1; j < len(G.History); j++ {
			if strings.HasPrefix(G.History[j], fenPrefix) {
				for k := j + 1; k < len(G.History); k++ {
					if strings.HasPrefix(G.History[k], fenPrefix) {
						return true
					}
				}
			}
		}
	}
	return false
}

func InsufficientMaterial(G *Game) bool {
	/*

		BUG!

		TODO:
		  	-(Any number of additional bishops of either color on the same color of square due to underpromotion do not affect the situation.)
	*/

	loneKing := []bool{
		G.Board.Occupied(WHITE)&G.Board.PieceBB[WHITE][KING] == G.Board.Occupied(WHITE),
		G.Board.Occupied(BLACK)&G.Board.PieceBB[BLACK][KING] == G.Board.Occupied(BLACK)}

	if !loneKing[WHITE] && !loneKing[BLACK] {
		return false
	}

	for color := WHITE; color <= BLACK; color++ {
		otherColor := []Color{BLACK, WHITE}[color]
		if loneKing[color] {
			// King vs King:
			if loneKing[otherColor] {
				return true
			}
			// King vs King & Knight
			if popcount(G.Board.PieceBB[otherColor][KNIGHT]) == 1 {
				mask := G.Board.PieceBB[otherColor][KING] | G.Board.PieceBB[otherColor][KNIGHT]
				occuppied := G.Board.Occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(G.Board.PieceBB[otherColor][BISHOP]) == 1 {
				mask := G.Board.PieceBB[otherColor][KING] | G.Board.PieceBB[otherColor][BISHOP]
				occuppied := G.Board.Occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := G.Board.PieceBB[color][KING] | G.Board.PieceBB[color][BISHOP]
		if (G.Board.Occupied(color)&kingBishopMask == G.Board.Occupied(color)) && (popcount(G.Board.PieceBB[color][BISHOP]) == 1) {
			mask := G.Board.PieceBB[otherColor][KING] | G.Board.PieceBB[otherColor][BISHOP]
			occuppied := G.Board.Occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(G.Board.PieceBB[otherColor][BISHOP]) == 1) {
				color1 := bitscan(G.Board.PieceBB[color][BISHOP]) % 2
				color2 := bitscan(G.Board.PieceBB[otherColor][BISHOP]) % 2
				if color1 == color2 {
					return true
				}
			}
		}

	}

	return false
}

/*******************************************************************************

	Functions that control the operating state of the game:

*******************************************************************************/

func (G *Game) GameOver(looser Color, reason string) {
	G.Result = []Color{BLACK, WHITE, DRAW}[looser] //oppoSite of the looser
	G.ResultDetail = reason
	G.Completed = true

	//TODO: Temporary:
	fmt.Println("{" + reason + "} " + []string{"1-0", "0-1", "1/2-1/2"}[G.Result])
}

/*******************************************************************************

	Modifiers:

*******************************************************************************/

func (G *Game) MakeMove(m Move) error {
	// TODO: 	Proper error checking along the way

	G.History = append(G.History, G.FEN())

	G.FiftyRule += 1

	G.MoveList = append(G.MoveList, m)

	from, to := getIndex(m.Algebraic)

	capturedColor, capturedPiece := G.Board.onSquare(to)

	if capturedPiece != NONE {
		// remove captured piece:
		G.Board.PieceBB[capturedColor][capturedPiece] ^= (1 << to)
		G.FiftyRule = 0
	}

	color, piece := G.Board.onSquare(from)
	if color == NEITHER || piece == NONE {
		return errors.New("Illegal Move.")
	}

	//move piece:
	G.Board.PieceBB[color][piece] ^= ((1 << from) | (1 << to))

	// Castle:
	if piece == KING {
		if from == (E1+56*uint8(color)) && (to == (G1 + 56*uint8(color))) {
			G.Board.PieceBB[color][ROOK] ^= (1 << (H1 + 56*uint8(color))) | (1 << (F1 + 56*uint8(color)))
		} else if from == (E1+56*uint8(color)) && to == (C1+56*uint8(color)) {
			G.Board.PieceBB[color][ROOK] ^= (1 << (A1 + 56*uint8(color))) | (1 << (D1 + 56*uint8(color)))
		}
	}

	// Cancel castling rights:
	for side := SHORT; side <= LONG; side++ {
		if piece == KING || //king moves
			(piece == ROOK && from == [2][2]uint8{{H1, A1}, {H8, A8}}[color][side]) {
			G.CastleRights[color][side] = false
		}
		if to == [2][2]uint8{{H8, A8}, {H1, A1}}[color][side] {
			G.CastleRights[[]Color{BLACK, WHITE}[color]][side] = false
		}
	}

	if piece == PAWN {
		G.FiftyRule = 0
		// Handle en Passant capture:
		if (G.EnPassant != 64) && (to == G.EnPassant) && (int(to)-int(from)%8 != 0) {
			if color == WHITE {
				G.Board.PieceBB[BLACK][PAWN] ^= (1 << (to - 8))
			} else {
				G.Board.PieceBB[WHITE][PAWN] ^= (1 << (to + 8))
			}
		}

		// Set en Passant:
		//if (((from / 8) + 1) == []uint8{2, 7}[color]) && (((to / 8) + 1) == []uint8{3, 6}[color]) {
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			G.EnPassant = uint8(int(from) + []int{8, -8}[color]) // type change fiasco! could crash.
		} else {
			G.EnPassant = 64
		}

		promotes := getPromotion(m.Algebraic)
		// Handle Promotions:
		if promotes != NONE {
			G.Board.PieceBB[color][piece] ^= (1 << to)    // remove pawn
			G.Board.PieceBB[color][promotes] ^= (1 << to) // add promoted piece
		}
	} else {
		G.EnPassant = 64
	}

	return nil
}

func (G *Game) FEN() string {

	piece := [][]string{
		{"P", "N", "B", "R", "Q", "K", " "},
		{"p", "n", "b", "r", "q", "k", " "},
		{" ", " ", " ", " ", " ", " ", " "}}

	var board string
	// put what is on each square into a squence (including blanks):
	for i := int(63); i >= 0; i-- {
		c, p := G.Board.onSquare(uint8(i))
		board += piece[c][p]
		if i%8 == 0 && i > 0 {
			board += "/"
		}
	}
	// replace groups of spaces with numbers instead
	for i := 8; i > 0; i-- {
		board = strings.Replace(board, strings.Repeat(" ", i), strconv.Itoa(i), -1)
	}
	// Player to move:
	turn := []string{"w", "b"}[G.toMove()]
	// Castling Rights:
	var rights string
	castles := [][]string{{"K", "Q"}, {"k", "q"}}
	for c := WHITE; c <= BLACK; c++ {
		for side := SHORT; side <= LONG; side++ {
			if G.CastleRights[c][side] {
				rights += castles[c][side]
			}
		}
	}
	if rights == "" {
		rights = "-"
	}
	// en Passant:
	var enPas string
	if G.EnPassant != 64 {
		enPas = getAlg(uint(G.EnPassant))
	} else {
		enPas = "-"
	}
	// Moves and 50 move rule
	fifty := strconv.Itoa(int(G.FiftyRule / 2))
	move := strconv.Itoa(int(len(G.MoveList)/2) + 1)
	// all together:
	fen := board + " " + turn + " " + rights + " " + enPas + " " + fifty + " " + move
	return fen
}

func (G *Game) LoadFEN(fen string) error {
	// TODO: error handling

	//root fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	words := strings.Split(fen, " ")

	// Move Count & toMove:
	unknownMove := Move{Algebraic: ""}
	fullMoves, _ := strconv.ParseUint(words[5], 10, 0)
	halfMoves := ((fullMoves - 1) * 2) + map[string]uint64{"w": 0, "b": 1}[words[1]]
	for i := uint64(0); i < halfMoves; i++ {
		G.MoveList = append(G.MoveList, unknownMove)
	}

	// 50 Move Rule:
	G.FiftyRule, _ = strconv.ParseUint(words[4], 10, 0)
	G.FiftyRule = (G.FiftyRule * 2) + map[string]uint64{"w": 0, "b": 1}[words[1]] //since internally we store half moves

	// Castling Rights:
	G.CastleRights = [2][2]bool{
		{strings.Contains(words[2], "K"), strings.Contains(words[2], "Q")},
		{strings.Contains(words[2], "k"), strings.Contains(words[2], "q")}}

	// en Passant:
	G.EnPassant = 64
	if words[3] != "-" {
		t, _ := strconv.ParseUint(words[3], 10, 0)
		G.EnPassant = uint8(t)
	}

	// Board position:
	piece := map[string]Piece{"P": PAWN, "p": PAWN,
		"N": KNIGHT, "n": KNIGHT,
		"B": BISHOP, "b": BISHOP,
		"R": ROOK, "r": ROOK,
		"Q": QUEEN, "q": QUEEN,
		"K": KING, "k": KING}
	color := map[string]Color{"P": WHITE, "p": BLACK,
		"N": WHITE, "n": BLACK,
		"B": WHITE, "b": BLACK,
		"R": WHITE, "r": BLACK,
		"Q": WHITE, "q": BLACK,
		"K": WHITE, "k": BLACK}

	G.Board.Clear()
	board := words[0]
	// remove the /'s and replace the numbers with that many spaces:
	parsedBoard := strings.Replace(board, "/", "", 9)
	for i := 1; i < 9; i++ {
		parsedBoard = strings.Replace(parsedBoard, strconv.Itoa(i), strings.Repeat(" ", i), -1)
	}
	// adjust the bitboards:
	for pos := 0; pos < len(parsedBoard); pos++ {
		k := parsedBoard[pos:(pos + 1)]
		_, ok := piece[k]
		if ok {
			G.Board.PieceBB[color[k]][piece[k]] |= (1 << uint(63-pos))
		}
	}
	return nil
}

func (G *Game) resetTimeControl() {
	G.MovesToGo = G.Moves
	G.Timer = [2]int64{G.Time, G.Time}
}

// TODO: change this to return a Game not act on one.
func (G *Game) initialize() error {
	// Sets up the game so that its ready for white to make the first move

	// TODO: this assumes a fresh unstarted game.
	G.resetTimeControl()
	//G.toMove = WHITE
	G.Board.Reset()
	G.CastleRights = [2][2]bool{{true, true}, {true, true}}
	G.EnPassant = 64
	G.Completed = false
	return nil
}

// Work in progress - should replace initialize():
func NewGame() Game {
	var g Game
	g.initialize()
	return g
}

/*******************************************************************************

	Getters:

*******************************************************************************/

func (G *Game) Print() {
	// TODO: this should instead take in *Game as an arguement
	G.Board.Print()
	fmt.Print([]string{"White", "Black"}[G.toMove()], " to move. ")
	if G.EnPassant != 64 {
		fmt.Print("Enpassant: ", getAlg(uint(G.EnPassant)), ". ")
	} else {
		fmt.Print("Enpassant: ", "none. ")
	}
	fmt.Print("Castling Rights: ",
		map[bool]string{true: "K", false: "-"}[G.CastleRights[WHITE][SHORT]],
		map[bool]string{true: "Q", false: "-"}[G.CastleRights[WHITE][LONG]],
		map[bool]string{true: "k", false: "-"}[G.CastleRights[BLACK][SHORT]],
		map[bool]string{true: "q", false: "-"}[G.CastleRights[BLACK][LONG]])
	fmt.Print("\n")
}

func (G *Game) PrintHUD() {
	toMove := G.toMove()
	lastMoveSource, lastMoveDestination := uint8(64), uint8(64)
	if len(G.MoveList) > 0 && G.MoveList[len(G.MoveList)-1].Algebraic != "" {
		lastMoveSource, lastMoveDestination = getIndex(G.MoveList[len(G.MoveList)-1].Algebraic)
	}
	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	fmt.Println("   +---+---+---+---+---+---+---+---+")
	for i := uint8(1); i <= 64; i++ {
		square := uint8(64 - i)
		if square%8 == 7 {
			fmt.Print(" ", square/8+1, " ")
		}
		fmt.Print("|")
		blankSquare := true
		for j := PAWN; j <= KING; j = j + 1 {
			for color := Color(0); color <= BLACK; color++ {
				if ((1 << square) & G.Board.PieceBB[color][j]) != 0 {
					if lastMoveDestination == square {
						fmt.Print("[", abbrev[color][j], "]")
					} else {
						fmt.Print(" ", abbrev[color][j], " ")
					}
					blankSquare = false
				}
			}
		}
		if blankSquare == true {
			if lastMoveSource == square {
				fmt.Print("[ ]")
			} else {
				fmt.Print("   ")
			}
		}
		if square%8 == 0 {
			fmt.Print("|")
			fmt.Print(strings.Repeat(" ", 6))
			switch square / 8 {
			case 7:
				formattedTimer := []string{"WHITE " + FormatTimer(G.Timer[WHITE]), "BLACK " + FormatTimer(G.Timer[BLACK])}
				formattedTimer[toMove] = "[" + formattedTimer[toMove] + "]"
				fmt.Print(formattedTimer[WHITE], strings.Repeat(" ", 3), formattedTimer[BLACK])
			case 6:
				fmt.Print("Move #: ", len(G.MoveList)/2, "    (Moves Remaining: ", G.MovesToGo, ")")

			case 5:
				if G.EnPassant != 64 {
					fmt.Print("Enpassant: ", getAlg(uint(G.EnPassant)))
				} else {
					fmt.Print("Enpassant: ", "None")
				}
			case 4:
				fmt.Print("Castling Rights: ",
					map[bool]string{true: "K", false: "-"}[G.CastleRights[WHITE][SHORT]],
					map[bool]string{true: "Q", false: "-"}[G.CastleRights[WHITE][LONG]],
					map[bool]string{true: "k", false: "-"}[G.CastleRights[BLACK][SHORT]],
					map[bool]string{true: "q", false: "-"}[G.CastleRights[BLACK][LONG]])
			case 3:
				fmt.Print("In Play: ", FormatPiecesInPlay(G)[WHITE])
			case 2:
				fmt.Print("         ", FormatPiecesInPlay(G)[BLACK])
			case 0:
				if len(G.MoveList) > 0 {
					fmt.Print("Last move: ", G.MoveList[len(G.MoveList)-1].Algebraic)
				}
			}

			fmt.Print("\n")
			fmt.Println("   +---+---+---+---+---+---+---+---+")
		}
	}
	fmt.Println("     a   b   c   d   e   f   g   h")
	title := G.Player[[]int{1, 0}[toMove]].Name + " (" + []string{"Black", "White"}[toMove] + ")"

	fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
	//var pv []string
	//if len(G.MoveList) > 0 {
	//	pv = G.MoveList[len(G.MoveList)-1].log
	//}
	//if len(pv)-2 >= 0 {
	//	fmt.Print(pv[len(pv)-2])
	//}
	title = G.Player[toMove].Name + " (" + []string{"White", "Black"}[toMove] + ")"
	fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
	/*
		fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
		pv = G.MoveList[len(G.MoveList)-1].log
		fmt.Print(pv[len(pv)-2])
		title = G.Player[toMove].Name + " (" + []string{"White", "Black"}[toMove] + ")"
		fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
	*/
}

func FormatPiecesInPlay(G *Game) [2]string {
	// Q RR BB NN PPPPPPPP
	// q rr bb nn pppppppp
	var pcs [2]string
	pieces := [][]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	counts := []int{8, 2, 2, 2, 1, 1}
	values := []int{1, 3, 3, 5, 9, 0}
	var scores [2]int
	for color := WHITE; color <= BLACK; color++ {
		for p := PAWN; p <= KING; p++ {
			inplay := int(popcount(G.Board.PieceBB[color][p]))
			dead := counts[p] - inplay
			scores[color] += (inplay * values[p])
			if dead < 0 {
				dead = 0
			}
			pcs[color] += strings.Repeat(pieces[color][p], inplay)
			pcs[color] += strings.Repeat(" ", dead+1)
		}
		pcs[color] += " (" + strconv.Itoa(scores[color]) + ")"
	}

	return pcs
}

func FormatTimer(ms int64) string {
	var r string
	if ms/60000 < 10 {
		r += "0"
	}
	r += strconv.FormatInt(ms/60000, 10) + ":"
	if (ms%60000)/1000 < 10 {
		r += "0"
	}
	r += strconv.FormatInt((ms%60000)/1000, 10) + "." + strconv.FormatInt((ms%60000)%1000, 10)
	if (ms%60000)%1000 < 10 {
		r += "00"
	} else if (ms%60000)%1000 < 100 {
		r += "0"
	}
	return r
}

func (G *Game) toMove() Color {
	return Color(len(G.MoveList) % 2)
}

func (G *Game) isInCheck(toMove Color) bool {
	// TODO: see isAttacked() notes
	notToMove := []Color{BLACK, WHITE}[toMove]
	kingsq := bitscan(G.Board.PieceBB[toMove][KING])
	return G.isAttacked(kingsq, notToMove)
}

func (G *Game) isAttacked(square uint, byWho Color) bool {
	// TODO: conceptually whether somebody is attacked or not isnt a property of the game
	//			but rather a property of the Player? So maybe have this be a stand alone function
	//			that takes in *Game
	defender := []Color{BLACK, WHITE}[byWho]

	// other king attacks:
	if (king_moves[square] & G.Board.PieceBB[byWho][KING]) != 0 {
		return true
	}

	// pawn attacks:
	if pawn_captures[defender][square]&G.Board.PieceBB[byWho][PAWN] != 0 {
		return true
	}

	// knight attacks:
	if knight_moves[square]&G.Board.PieceBB[byWho][KNIGHT] != 0 {
		return true
	}
	// diagonal attacks:
	direction := [4][65]uint64{nw, ne, sw, se}
	scan := [4]func(uint64) uint{BSF, BSF, BSR, BSR}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.Board.Occupied(BOTH))
		if (1<<blockerIndex)&(G.Board.PieceBB[byWho][BISHOP]|G.Board.PieceBB[byWho][QUEEN]) != 0 {
			return true
		}
	}
	// straight attacks:
	direction = [4][65]uint64{north, west, south, east}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.Board.Occupied(BOTH))
		if (1<<blockerIndex)&(G.Board.PieceBB[byWho][ROOK]|G.Board.PieceBB[byWho][QUEEN]) != 0 {
			return true
		}
	}
	return false
}

/*******************************************************************************

	Game Log:

*******************************************************************************/

func (G *Game) StartLog() error {
	fmt.Print("Creating log file... ")

	//check if folder exists:
	//if err := os.Mkdir("logs", os.ModePerm); !os.IsExist(err) {
	//	return err
	//}
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	//check if the file exists:
	filename := fmt.Sprint("logs/", G.Event, " round ", G.Round, ".log")
	if _, test := os.Stat(filename); os.IsNotExist(test) {
		// file doesnt exist
	} else if test == nil {
		// file does exist
		os.Remove(filename)
	}
	var err error
	G.logFile, err = os.Create(filename)
	if err != nil {
		return err
	}

	// Give header information for the log file:
	G.logBuffer += fmt.Sprintln(G.Event) +
		fmt.Sprintln("Round ", G.Round) +
		fmt.Sprintln(G.Player[WHITE].Name, "vs", G.Player[BLACK].Name) +
		fmt.Sprintln(time.Now().Format("01/02/2006 15:04:05.000")) +
		fmt.Sprintln("")
	err = G.AppendLog()

	fmt.Println("Success.")
	return err
}

func (G *Game) AppendLog() error {
	if _, err := G.logFile.WriteString(G.logBuffer); err != nil {
		return err
	}
	G.logBuffer = ""
	return nil
}

func (G *Game) CloseLog() error {
	return G.logFile.Close()
}
