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


*******************************************************************************/

package main

import (
	"fmt"
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
	// TODO: 50 move rule and 3 fold

	// Header info (General info, usually needed for pgn):
	time      int64
	moves     int64
	repeating bool
	player    [2]Engine //white=0,black=1

	//TODO: the rest of the PGN required heading. i think there are 8?

	// Control info:
	timer     [2]int64 //the game clock for each side. milliseconds.
	movesToGo int64    //moves until time control

	// Game run time info:
	state        Status //UNSTARTED,RUNNING,STOPPED
	board        Board
	fiftyRule    uint64
	history      []string // FEN of every known position so far.
	enPassant    uint8
	castleRights [2][2]bool
	moveList     []Move //move history
	completed    bool   // TODO:  change this to reflect how the game ended. time, checkmate, adjunction, etc

	// Post game info:
	result Color //WHITE,BLACK,DRAW - should be set when game state is changed to STOPPED

}

/*******************************************************************************

	Primary function. Manages the game itself. Loops until instructed not to.
	Or until the game ends.

*******************************************************************************/

func (G *Game) playLoop() error {
	// TODO:
	//		-Optimization: legal move gen gets called twice per loop. once in the notation adaptation, and once in move verification
	for G.state == RUNNING {
		for color := WHITE; color <= BLACK; color++ {
			otherColor := []Color{BLACK, WHITE}[color]
			// Tell the engine to set its internal board:
			if err := G.player[color].Set(G.moveList); err != nil { //TODO: should probably pass by ref
				G.Stop()
				break
			}
			// Request a move from the engine:
			startTime := time.Now()
			move, err := G.player[color].Move(G.timer, G.movesToGo)
			if err != nil {
				G.Stop()
				break
			}
			endTime := time.Now()
			lapsed := endTime.Sub(startTime)
			// Adjust time control:
			G.timer[color] -= int64(lapsed.Seconds() * 1000)
			if G.timer[color] <= 0 {
				G.GameOver(color, "Out of time. Used "+strconv.FormatInt(int64(lapsed.Seconds()*1000), 10)+" ms.")
				break
			}
			// Convert the notation from the engines notation to pure coordinate notation
			preparsedMove := move
			move.algebraic = InternalizeNotation(G, preparsedMove.algebraic)

			fmt.Print([]string{"WHITE", "BLACK"}[color], "> ", move.algebraic, " (From Engine: ", preparsedMove.algebraic, ")\n")
			// Check legality of move.
			LegalMoves := LegalMoveList(G)
			if !contains(LegalMoves, move) {
				G.GameOver(color, "Illegal move.")
				break
			}
			// Adjust the internal board:
			if err = G.MakeMove(move); err != nil {
				G.GameOver(color, err.Error()) // illegal move
				break
			}
			// Check:
			check := G.isInCheck(otherColor)
			// Mate:
			oppLegalMoves := LegalMoveList(G)
			if len(oppLegalMoves) == 0 {
				if check {
					//checkmate!
					G.GameOver(otherColor, "Checkmate.")
					break
				} else {
					//stalemate!
					G.GameOver(NEITHER, "Stalemate.")
					break
				}
			}
			// 50 Move Draw:
			if FiftyMoveDraw(G) {
				G.GameOver(NEITHER, "50 move rule.")
				break
			}
			// Insufficient material:
			if InsufficientMaterial(G) {
				G.GameOver(NEITHER, "Insufficient material.")
				break
			}
			// 3 fold:
			if ThreeFold(G) {
				G.GameOver(NEITHER, "Three fold repitition.")
				break
			}
		}
		if G.completed == false {
			G.movesToGo -= 1
			if G.movesToGo == 0 && G.repeating == true {
				G.resetTimeControl()
			}
		}
	}
	G.board.Print()
	return nil
}

func contains(list []Move, move Move) bool {
	for i, _ := range list {
		if move.algebraic == list[i].algebraic {
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
	// maybe sort them alphabetically first
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
		// King vs King & opposite bishop
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

func (G *Game) Start() error {
	// TODO: Start the game with error checking based on G.state and G.completed
	if G.state == UNSTARTED {
		G.initialize()
	}
	G.state = RUNNING
	fmt.Println("Game running. (", G.player[WHITE].Name, "vs", G.player[BLACK].Name, ")") //TODO: include round#

	// Start up the engines:
	G.player[WHITE].Start()
	G.player[BLACK].Start()

	// Begin playing the game:
	G.playLoop()

	return nil
}

func (G *Game) Stop() error {
	// TODO: Stop the game with error checking based on G.state and G.completed
	if G.state == RUNNING {
		// Turn off the engines:
		G.player[WHITE].Shutdown()
		G.player[BLACK].Shutdown()

		G.state = STOPPED
		fmt.Println("Game stopped.")
	}
	return nil
}

func (G *Game) GameOver(looser Color, reason string) {
	fmt.Println("Game Over.", []string{"White looses.", "Black looses.", "Draw."}[looser], reason)
	G.result = []Color{BLACK, WHITE, DRAW}[looser] //opposite of the looser
	G.completed = true
	G.Stop()
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

	G.moveList = append(G.moveList, m)

	from, to := getIndex(m.algebraic)

	capturedColor, capturedPiece := G.board.onSquare(to)

	if capturedPiece != NONE {
		// remove captured piece:
		G.board.pieceBB[capturedColor][capturedPiece] ^= (1 << to)
		G.fiftyRule = 0
	}

	color, piece := G.board.onSquare(from)
	if color == NEITHER || piece == NONE {
		return customError{"Illegal Move", time.Now()}
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

		promotes := getPromotion(m.algebraic)
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
	move := strconv.Itoa(int(len(G.moveList)/2) + 1)
	// all together:
	fen := board + " " + turn + " " + rights + " " + enPas + " " + fifty + " " + move
	return fen
}

func (G *Game) LoadFEN(fen string) error {
	// TODO: error handling

	//root fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	words := strings.Split(fen, " ")

	// Move Count & toMove:
	unknownMove := Move{algebraic: ""}
	fullMoves, _ := strconv.ParseUint(words[5], 10, 0)
	halfMoves := ((fullMoves - 1) * 2) + map[string]uint64{"w": 0, "b": 1}[words[1]]
	for i := uint64(0); i < halfMoves; i++ {
		G.moveList = append(G.moveList, unknownMove)
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
	G.state = UNSTARTED
	G.completed = false
	return nil
}

func (G *Game) Validate() error {
	// TODO: Should check that all the data members are set up correctly to not cause a crash.
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

func (G *Game) toMove() Color {
	return Color(len(G.moveList) % 2)
}

func (G *Game) isInCheck(toMove Color) bool {
	// TODO: see isAttacked() notes
	notToMove := []Color{BLACK, WHITE}[toMove]
	kingsq := bitscan(G.board.pieceBB[toMove][KING])
	return G.isAttacked(kingsq, notToMove)
}

func (G *Game) isAttacked(square uint, byWho Color) bool {
	// TODO: conceptually whether somebody is attacked or not isnt a property of the game
	//			but rather a property of the player? So maybe have this be a stand alone function
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
