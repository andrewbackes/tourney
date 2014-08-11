/*******************************************************************************

 Project: Tourney

 Module: notation
 Description: notation parsing utilities

 The internals of Tourney use Pure Coordinate Notation to represent moves. For
 example, a2a4 means to move what is on a2 to a4. However, UCI and WB require
 Standard Algebraic Notation. The purpose of this module is to translate between
 PCN and SAN.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/29/2014


 TODO:
 	-instead of calling LegalMoveGen(), should have a stripped down version
 	for just the type of piece moving.

*******************************************************************************/

/*
BUG:


> ucinewgame
> position startpos moves d2d4 g8f6 c2c4 e7e6 b1c3 f8b4 g1f3 c7c5 d4c5 b8c6 c1d2 b4c5 e2e3 d7d5 f1e2 e8g8 e1g1 d5c4 e2c4 e6e5 e3e4 c8g4 a1c1 c6d4 c4e2 d4f3 e2f3 g4f3 g2f3 c5d4 d1e2 d8b6 d2e3 a8c8 f1d1 f8d8 g1g2 a7a6 d1d3 f6d7 c3d5 b6g6 g2h1 c8c1 e3c1 g6e6 c1g5 f7f6 g5e3 d7b6 e3d4 e5d4 e2d1 b6d5 d3d4 g8f7 d4d5 d8d5 e4d5 e6d6 h1g2 f6f5 d1d3 f7g6 h2h4 g6f6 d3d4 f6g6 f3f4 d6e7 d5d6 e7d7 g2g3 h7h6 d4c5 d7e6 c5e5 e6f7 b2b3 g6h7 g3f3 h7g6 f3e3 b7b5 e3d4 f7a7 e5c5 a7b7 h4h5 g6f6 c5e5 f6f7 e5f5 f7e8 f5e6 e8f8 d6d7 b7b8 d4c5 b8c7 c5d5 c7d8 d5c6 d8a8 c6b6 a8d8 b6a6 d8b8 e6e5 b8d8 e5d6 f8f7 a6b5 f7g8 f4f5 g8h8 a2a4 h8g8 b3b4 g7g6 f5g6 g8g7 d6d4 g7f8 b5c6 f8g8 d4d5 g8h8 d5e5 h8g8 e5e8 d8e8 d7e8Q g8g7 e8f7 g7h8 f7f8
Error: index out of range.
5Q1k/8/2K3Pp/7P/PP6/8/5P2/8 b - - 2 69
0000
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |[Q]|   | k |      WHITE 00:00.285   [BLACK 00:00.970]
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |[ ]|   |   |
+---+---+---+---+---+---+---+---+
|   |   | K |   |   |   | P | p |      Enpassant: None
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   | P |      Castling Rights: ----
+---+---+---+---+---+---+---+---+
| P | P |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   | P |   |   |
+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   |   |   |
+---+---+---+---+---+---+---+---+
panic: runtime error: index out of range

goroutine 21 [running]:
runtime.panic(0x15ace0, 0x2551fc)
	/usr/local/go/src/pkg/runtime/panic.c:279 +0xf5
main.InternalizeNotation(0x2082fc1c0, 0x208462bc9, 0x4, 0x0, 0x0)
	/Users/Andrew/Projects/Tourney/src/tourney/notation.go:72 +0xb2f
main.ExecuteNextTurn(0x2082fc1c0, 0x0)
	/Users/Andrew/Projects/Tourney/src/tourney/game.go:147 +0x561
main.PlayGame(0x2082fc1c0, 0x0, 0x0)
	/Users/Andrew/Projects/Tourney/src/tourney/game.go:100 +0x12d
main.RunTourney(0x2082c8dd0, 0x0, 0x0)
	/Users/Andrew/Projects/Tourney/src/tourney/tourney.go:131 +0xfc5
main.func·004()
	/Users/Andrew/Projects/Tourney/src/tourney/command.go:42 +0x32
created by main.doCommand
	/Users/Andrew/Projects/Tourney/src/tourney/command.go:46 +0x689

*/

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func InternalizeNotation(G *Game, moveToParse string) string {
	// Converts Standard Algebraic Notation (SAN) to Pure Coordinate Notation (PCN)
	// Examples of PCN are: e2e4 (and) e7e8Q

	// TODO: needs error handling.
	// TODO: what about promotion captures? or ambiguous promotions?

	// First check to see if it is already in the correct form.
	PCN := "([a-h][1-8])([a-h][1-8])([QBNRqbnr]?)"
	matches, _ := regexp.MatchString(PCN, moveToParse)
	if matches {
		parsed := moveToParse[:len(moveToParse)-1]
		// Some engines dont capitalize the promotion piece:
		parsed += strings.ToUpper(moveToParse[len(moveToParse)-1:])
		return parsed
	}
	// Check for null move:
	if moveToParse == "0000" {
		return moveToParse
	}

	// Check for castling:
	if moveToParse == "O-O" {
		return []string{"e1g1", "e8g8"}[G.toMove()]
	}
	if moveToParse == "O-O-O" {
		return []string{"e1c1", "e8c8"}[G.toMove()]
	}

	// First check for an ambiguous promotion:
	//SAN := "([a-h]?)([0-9]?)([a-h][0-9])([=])([BNQR])([+#]?)"
	//p, _ := regexp.Compile(SAN)

	// Breakdown the SAN:
	SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x=]?)([BKNPQR]|[a-h][1-8])([+#]?)"
	r, _ := regexp.Compile(SAN)

	matched := r.FindStringSubmatch(moveToParse)
	if len(matched) == 0 {
		fmt.Println("Error: index out of range.")
		fmt.Println(G.FEN())
		fmt.Println(moveToParse)
		G.PrintHUD()
	}
	// For the sake of sanity, lets name some stuff:
	piece := matched[1]
	fromFile := matched[2]
	fromRank := matched[3]
	action := matched[4]      // capture or promote
	destination := matched[5] //or promotion piece if action="="
	//check := matched[6]       //or mate
	var promote string

	if piece == "" {
		piece = "P"
	}

	// TODO: what about promotion captures? or ambiguous promotions?
	if action == "=" {
		promote = destination
		destination = fromFile + fromRank
		fromFile = ""
		fromRank = ""
	}

	origin, err := originOfPiece(piece, destination, fromFile, fromRank, G)
	if err != nil {
		fmt.Println(err)
		fmt.Println(G.FEN())
		fmt.Println(moveToParse)
		G.PrintHUD()
	}

	return origin + destination + promote
}

func originOfPiece(piece, destination, fromFile, fromRank string, G *Game) (string, error) {
	pieceMap := map[string]Piece{
		"P": PAWN, "p": PAWN,
		"N": KNIGHT, "n": KNIGHT,
		"B": BISHOP, "b": BISHOP,
		"R": ROOK, "r": ROOK,
		"Q": QUEEN, "q": QUEEN,
		"K": KING, "k": KING}

	if fromFile != "" && fromRank != "" {
		return fromFile + fromRank, nil
	}

	// Get all legal moves:
	legalMoves := LegalMoveList(G)
	var eligableMoves []Move

	// Grab the legal moves that land on our square:
	for _, mv := range legalMoves {
		dest := mv.Algebraic[2:4]
		if dest == destination {
			eligableMoves = append(eligableMoves, mv)
		}
	}

	// Get all the squares that have our piece on it from the move list:
	color := G.toMove()
	var eligableSquares []string
	bits := G.board.pieceBB[color][pieceMap[piece]]
	for bits != 0 {
		bit := bitscan(bits)
		sq := getAlg(bit)
		//verify that its a legal move:
		for _, mv := range eligableMoves {
			if mv.Algebraic[:2] == sq {
				eligableSquares = append(eligableSquares, sq)
				break
			}
		}
		bits ^= (1 << bit)
	}

	// Look for one of the squares that matches the file/rank criteria:
	for _, sq := range eligableSquares {
		if ((sq[0:1] == fromFile) || (fromFile == "")) && ((sq[1:2] == fromRank) || (fromRank == "")) {
			return sq, nil
		}

	}
	return "", errors.New("Notation: Can not find source square.")
}

/*******************************************************************************

	Notation Stuff:

*******************************************************************************/

func StripAnnotations(mv string) string {
	m := strings.Replace(mv, "!", "", -1)
	m = strings.Replace(m, "?", "", -1)
	return m
}

func getIndex(alg string) (uint8, uint8) {
	// TODO: accept more notation.

	// For the form: e2e4
	from := []byte{alg[0], alg[1]}
	to := []byte{alg[2], alg[3]}

	f := ((from[1] - 48) * 8) - (from[0] - 96)
	t := ((to[1] - 48) * 8) - (to[0] - 96)

	// TODO: of the form e2-e4

	// TODO: of the form Nc3

	return f, t
}

func getPromotion(alg string) Piece {
	// TODO: accept more notation. Currently only accepts e7e8Q, e7e8=Q, and e7e8/Q
	if len(alg) > 4 {
		p := make(map[string]Piece)
		p = map[string]Piece{"Q": QUEEN, "N": KNIGHT, "B": BISHOP, "R": ROOK}
		return p[string(alg[len(alg)-1])]
	}
	return NONE
}

func getAlg(index uint) string {
	var r string

	file := rune(97 + (7 - (index % 8)))
	rank := rune((index / 8) + 49)

	r = string(file) + string(rank)

	return r
}

/*******************************************************************************

	Tests:

*******************************************************************************/

// TODO: add to a _test file

/*

func TestNotation() {
	var G Game
	G.board.Reset()

	G.LoadFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	G.Print()
	moves := []string{"Qh3", "Nb5", "Pxh3", "xh3", "O-O-O", "O-O"}
	for _, m := range moves {
		fmt.Print(m, "->")
		n := InternalizeNotation(&G, m)
		fmt.Print(n, "\n")
	}

	G.LoadFEN("8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1")
	G.Print()
	moves = []string{"gf4", "ef4", "gxf4", "Rf4", "Rxf4"}
	for _, m := range moves {
		fmt.Print(m, "->")
		n := InternalizeNotation(&G, m)
		fmt.Print(n, "\n")
	}

	return
}

*/
