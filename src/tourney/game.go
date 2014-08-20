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


*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	//"runtime"
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
	time        int64
	moves       int64
	repeating   bool
	Player      [2]Engine // white=0,black=1
	StartingFEN string    // first position out of the book

	// Control info:
	timer     [2]int64 //the game clock for each side. milliseconds.
	movesToGo int64    //moves until time control

	// Game run time info:
	board        Board
	fiftyRule    uint64
	history      []string // FEN of every known position so far.
	enPassant    uint8
	castleRights [2][2]bool
	MoveList     []Move //move history
	Completed    bool   // TODO:  change this to reflect how the game ended. time, checkmate, adjunction, etc

	// Post game info:
	Result       Color //WHITE,BLACK,DRAW
	ResultDetail string
}

/*******************************************************************************

	Primary function. Manages the game itself. Loops until instructed not to.
	Or until the game ends.

*******************************************************************************/

func PlayGame(G *Game) error {
	// Note: opening book is handled in RunTourney()
	// Start up the engines:
	if err := G.Player[WHITE].Start(); err != nil {
		return err
	}
	if err := G.Player[BLACK].Start(); err != nil {
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
		if gameover {
			state = STOPPED
			break
		}

		if G.toMove() == WHITE {
			G.movesToGo -= 1
			if G.movesToGo == 0 && G.repeating == true {
				G.resetTimeControl()
			}
		}
	}

	// Stop the engines:
	G.Player[WHITE].Shutdown()
	G.Player[BLACK].Shutdown()

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
	startTime := time.Now() // TODO: move the timing stuff to the Move() method.
	move, err := G.Player[color].Move(G.timer, G.movesToGo)
	if err != nil {
		G.GameOver(color, err.Error())
		return true
	}
	endTime := time.Now()
	lapsed := endTime.Sub(startTime)
	// Adjust time control:
	G.timer[color] -= int64(lapsed.Seconds() * 1000)
	if G.timer[color] <= 0 {
		G.GameOver(color, "Out of time. Used "+strconv.FormatInt(int64(lapsed.Seconds()*1000), 10)+" ms.")
		return true
	}
	// Convert the notation from the engines notation to pure coordinate notation
	preparsedMove := move
	move.Algebraic = InternalizeNotation(G, preparsedMove.Algebraic)

	fmt.Print([]string{"WHITE", "BLACK"}[color], "> ", move.Algebraic, " (From Engine: ", preparsedMove.Algebraic, ")\n")
	// Check legality of move.
	LegalMoves := LegalMoveList(G)
	if !contains(LegalMoves, move) {
		G.GameOver(color, "Illegal move.")
		return true
	}
	// Adjust the internal board:
	if err = G.MakeMove(move); err != nil {
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
	return G.fiftyRule == 100
}

func ThreeFold(G *Game) bool {
	// TODO: Can this be optimized? very very slow right now.
	// maybe sort them alphabetically first.
	// should go backwards instead, as many times as in the fiftyMove count.
	for i := 0; i < len(G.history); i++ {
		fen := G.history[i]
		fenSplit := strings.Split(fen, " ")
		fenPrefix := fenSplit[0] + " " + fenSplit[1] + " " + fenSplit[2] + " " + fenSplit[3]
		for j := i + 1; j < len(G.history); j++ {
			if strings.HasPrefix(G.history[j], fenPrefix) {
				for k := j + 1; k < len(G.history); k++ {
					if strings.HasPrefix(G.history[k], fenPrefix) {
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
		G.board.Occupied(WHITE)&G.board.pieceBB[WHITE][KING] == G.board.Occupied(WHITE),
		G.board.Occupied(BLACK)&G.board.pieceBB[BLACK][KING] == G.board.Occupied(BLACK)}

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
			if popcount(G.board.pieceBB[otherColor][KNIGHT]) == 1 {
				mask := G.board.pieceBB[otherColor][KING] | G.board.pieceBB[otherColor][KNIGHT]
				occuppied := G.board.Occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
			// King vs King & Bishop
			if popcount(G.board.pieceBB[otherColor][BISHOP]) == 1 {
				mask := G.board.pieceBB[otherColor][KING] | G.board.pieceBB[otherColor][BISHOP]
				occuppied := G.board.Occupied(otherColor)
				if occuppied&mask == occuppied {
					return true
				}
			}
		}
		// King vs King & oppoSite bishop
		kingBishopMask := G.board.pieceBB[color][KING] | G.board.pieceBB[color][BISHOP]
		if (G.board.Occupied(color)&kingBishopMask == G.board.Occupied(color)) && (popcount(G.board.pieceBB[color][BISHOP]) == 1) {
			mask := G.board.pieceBB[otherColor][KING] | G.board.pieceBB[otherColor][BISHOP]
			occuppied := G.board.Occupied(otherColor)
			if (occuppied&mask == occuppied) && (popcount(G.board.pieceBB[otherColor][BISHOP]) == 1) {
				color1 := bitscan(G.board.pieceBB[color][BISHOP]) % 2
				color2 := bitscan(G.board.pieceBB[otherColor][BISHOP]) % 2
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
	fmt.Println("Game Over.", []string{"White looses.", "Black looses.", "Draw."}[looser], reason)
	G.Result = []Color{BLACK, WHITE, DRAW}[looser] //oppoSite of the looser
	G.ResultDetail = reason
	G.Completed = true
}

/*******************************************************************************

	Modifiers:

*******************************************************************************/

func (G *Game) MakeMove(m Move) error {
	// TODO:	Other notation
	// TODO: 	Proper error checking along the way
	//			promotions

	G.history = append(G.history, G.FEN())

	G.fiftyRule += 1

	G.MoveList = append(G.MoveList, m)

	from, to := getIndex(m.Algebraic)

	capturedColor, capturedPiece := G.board.onSquare(to)

	if capturedPiece != NONE {
		// remove captured piece:
		G.board.pieceBB[capturedColor][capturedPiece] ^= (1 << to)
		G.fiftyRule = 0
	}

	color, piece := G.board.onSquare(from)
	if color == NEITHER || piece == NONE {
		return errors.New("Illegal Move.")
	}

	//move piece:
	G.board.pieceBB[color][piece] ^= ((1 << from) | (1 << to))

	// Castle:
	if piece == KING {
		if from == (E1+56*uint8(color)) && (to == (G1 + 56*uint8(color))) {
			G.board.pieceBB[color][ROOK] ^= (1 << (H1 + 56*uint8(color))) | (1 << (F1 + 56*uint8(color)))
		} else if from == (E1+56*uint8(color)) && to == (C1+56*uint8(color)) {
			G.board.pieceBB[color][ROOK] ^= (1 << (A1 + 56*uint8(color))) | (1 << (D1 + 56*uint8(color)))
		}
	}

	// Cancel castling rights:
	for side := SHORT; side <= LONG; side++ {
		if piece == KING || //king moves
			(piece == ROOK && from == [2][2]uint8{{H1, A1}, {H8, A8}}[color][side]) {
			G.castleRights[color][side] = false
		}
		if to == [2][2]uint8{{H8, A8}, {H1, A1}}[color][side] {
			G.castleRights[[]Color{BLACK, WHITE}[color]][side] = false
		}
	}

	if piece == PAWN {
		G.fiftyRule = 0
		// Handle en Passant capture:
		if (G.enPassant != 64) && (to == G.enPassant) && (int(to)-int(from)%8 != 0) {
			if color == WHITE {
				G.board.pieceBB[BLACK][PAWN] ^= (1 << (to - 8))
			} else {
				G.board.pieceBB[WHITE][PAWN] ^= (1 << (to + 8))
			}
		}

		// Set en Passant:
		//if (((from / 8) + 1) == []uint8{2, 7}[color]) && (((to / 8) + 1) == []uint8{3, 6}[color]) {
		if int(from)-int(to) == 16 || int(from)-int(to) == -16 {
			G.enPassant = uint8(int(from) + []int{8, -8}[color]) // type change fiasco! could crash.
		} else {
			G.enPassant = 64
		}

		promotes := getPromotion(m.Algebraic)
		// Handle Promotions:
		if promotes != NONE {
			G.board.pieceBB[color][piece] ^= (1 << to)    // remove pawn
			G.board.pieceBB[color][promotes] ^= (1 << to) // add promoted piece
		}
	} else {
		G.enPassant = 64
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
		c, p := G.board.onSquare(uint8(i))
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
			if G.castleRights[c][side] {
				rights += castles[c][side]
			}
		}
	}
	if rights == "" {
		rights = "-"
	}
	// en Passant:
	var enPas string
	if G.enPassant != 64 {
		enPas = getAlg(uint(G.enPassant))
	} else {
		enPas = "-"
	}
	// Moves and 50 move rule
	fifty := strconv.Itoa(int(G.fiftyRule / 2))
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
	G.fiftyRule, _ = strconv.ParseUint(words[4], 10, 0)
	G.fiftyRule = (G.fiftyRule * 2) + map[string]uint64{"w": 0, "b": 1}[words[1]] //since internally we store half moves

	// Castling Rights:
	G.castleRights = [2][2]bool{
		{strings.Contains(words[2], "K"), strings.Contains(words[2], "Q")},
		{strings.Contains(words[2], "k"), strings.Contains(words[2], "q")}}

	// en Passant:
	G.enPassant = 64
	if words[3] != "-" {
		t, _ := strconv.ParseUint(words[3], 10, 0)
		G.enPassant = uint8(t)
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

	G.board.Clear()
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
			G.board.pieceBB[color[k]][piece[k]] |= (1 << uint(63-pos))
		}
	}
	return nil
}

func (G *Game) resetTimeControl() {
	G.movesToGo = G.moves
	G.timer = [2]int64{G.time, G.time}
}

func (G *Game) initialize() error {
	// Sets up the game so that its ready for white to make the first move

	// TODO: this assumes a fresh unstarted game.
	G.resetTimeControl()
	//G.toMove = WHITE
	G.board.Reset()
	G.castleRights = [2][2]bool{{true, true}, {true, true}}
	G.enPassant = 64
	G.Completed = false
	return nil
}

/*******************************************************************************

	Getters:

*******************************************************************************/

func (G *Game) Print() {
	// TODO: this should instead take in *Game as an arguement
	G.board.Print()
	fmt.Print([]string{"White", "Black"}[G.toMove()], " to move. ")
	if G.enPassant != 64 {
		fmt.Print("Enpassant: ", getAlg(uint(G.enPassant)), ". ")
	} else {
		fmt.Print("Enpassant: ", "none. ")
	}
	fmt.Print("Castling Rights: ",
		map[bool]string{true: "K", false: "-"}[G.castleRights[WHITE][SHORT]],
		map[bool]string{true: "Q", false: "-"}[G.castleRights[WHITE][LONG]],
		map[bool]string{true: "k", false: "-"}[G.castleRights[BLACK][SHORT]],
		map[bool]string{true: "q", false: "-"}[G.castleRights[BLACK][LONG]])
	fmt.Print("\n")
}

func (G *Game) PrintHUD() {
	toMove := G.toMove()
	lastMoveSource, lastMoveDestination := uint8(64), uint8(64)
	if len(G.MoveList) > 0 {
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
				if ((1 << square) & G.board.pieceBB[color][j]) != 0 {
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
				formattedTimer := []string{"WHITE " + FormatTimer(G.timer[WHITE]), "BLACK " + FormatTimer(G.timer[BLACK])}
				formattedTimer[toMove] = "[" + formattedTimer[toMove] + "]"
				fmt.Print(formattedTimer[WHITE], strings.Repeat(" ", 3), formattedTimer[BLACK])
			case 6:
				fmt.Print("Move #: ", len(G.MoveList)/2, "    (Moves Remaining: ", G.movesToGo, ")")

			case 5:
				if G.enPassant != 64 {
					fmt.Print("Enpassant: ", getAlg(uint(G.enPassant)))
				} else {
					fmt.Print("Enpassant: ", "None")
				}
			case 4:
				fmt.Print("Castling Rights: ",
					map[bool]string{true: "K", false: "-"}[G.castleRights[WHITE][SHORT]],
					map[bool]string{true: "Q", false: "-"}[G.castleRights[WHITE][LONG]],
					map[bool]string{true: "k", false: "-"}[G.castleRights[BLACK][SHORT]],
					map[bool]string{true: "q", false: "-"}[G.castleRights[BLACK][LONG]])
			case 3:
				fmt.Print("Out of Play: ", FormatGraveyard(G)[WHITE])
			case 2:
				fmt.Print("             ", FormatGraveyard(G)[BLACK])
			case 0:
				fmt.Print("Last move: ", G.MoveList[len(G.MoveList)-1].Algebraic)
			}

			fmt.Print("\n")
			fmt.Println("   +---+---+---+---+---+---+---+---+")
		}
	}
	fmt.Println("     a   b   c   d   e   f   g   h")
	title := G.Player[[]int{1, 0}[toMove]].Name + " (" + []string{"Black", "White"}[toMove] + ")"

	fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
	pv := G.MoveList[len(G.MoveList)-1].log
	fmt.Print(pv[len(pv)-2])
	title = G.Player[toMove].Name + " (" + []string{"White", "Black"}[toMove] + ")"
	fmt.Print(strings.Repeat("-", (80-len(title))/2), title, strings.Repeat("-", (80-len(title))/2), "\n")
}

func FormatGraveyard(G *Game) [2]string {
	// Q RR BB NN PPPPPPPP
	// q rr bb nn pppppppp
	var gy [2]string
	pieces := [][]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	counts := []int{8, 2, 2, 2, 1, 1}
	for color := WHITE; color <= BLACK; color++ {
		for p := PAWN; p < KING; p++ {
			inplay := int(popcount(G.board.pieceBB[color][p]))
			dead := counts[p] - inplay
			if dead < 0 {
				dead = 0
			}
			gy[color] += strings.Repeat(pieces[color][p], dead)
			gy[color] += strings.Repeat(" ", inplay+1)
		}
	}

	return gy
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
	kingsq := bitscan(G.board.pieceBB[toMove][KING])
	return G.isAttacked(kingsq, notToMove)
}

func (G *Game) isAttacked(square uint, byWho Color) bool {
	// TODO: conceptually whether somebody is attacked or not isnt a property of the game
	//			but rather a property of the Player? So maybe have this be a stand alone function
	//			that takes in *Game
	defender := []Color{BLACK, WHITE}[byWho]

	// other king attacks:
	if (king_moves[square] & G.board.pieceBB[byWho][KING]) != 0 {
		return true
	}

	// pawn attacks:
	if pawn_captures[defender][square]&G.board.pieceBB[byWho][PAWN] != 0 {
		return true
	}

	// knight attacks:
	if knight_moves[square]&G.board.pieceBB[byWho][KNIGHT] != 0 {
		return true
	}
	// diagonal attacks:
	direction := [4][65]uint64{nw, ne, sw, se}
	scan := [4]func(uint64) uint{BSF, BSF, BSR, BSR}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.board.Occupied(BOTH))
		if (1<<blockerIndex)&(G.board.pieceBB[byWho][BISHOP]|G.board.pieceBB[byWho][QUEEN]) != 0 {
			return true
		}
	}
	// straight attacks:
	direction = [4][65]uint64{north, west, south, east}
	for i := 0; i < 4; i++ {
		blockerIndex := scan[i](direction[i][square] & G.board.Occupied(BOTH))
		if (1<<blockerIndex)&(G.board.pieceBB[byWho][ROOK]|G.board.pieceBB[byWho][QUEEN]) != 0 {
			return true
		}
	}
	return false
}
