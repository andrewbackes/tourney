/*******************************************************************************

 Project: Tourney

 Module: notation
 Description: notation parsing utilities

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/29/2014

*******************************************************************************/

package main

import (
	//"fmt"
	"regexp"
)

func InternalizeNotation(G *Game, moveToParse string) string {
	// Converts Standard Algebraic Notation (SAN) to Pure Coordinate Notation (PCN)
	// Examples of PCN are: e2e4 (and) e7e8Q

	// TODO: what about promotion captures? or ambiguous promotions?

	// First check to see if it is already in the correct form.
	PCN := "([a-h][1-8])([a-h][1-8])([QBNR]?)"
	matches, _ := regexp.MatchString(PCN, moveToParse)
	if matches {
		return moveToParse
	}
	// Check for castling:
	if moveToParse == "O-O" {
		return []string{"e1g1", "e8g8"}[G.toMove()]
	}
	if moveToParse == "O-O-O" {
		return []string{"e1c1", "e8c8"}[G.toMove()]
	}

	// Breakdown the SAN:
	SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x=]?)([BKNPQR]|[a-h][1-8])([+#]?)"
	r, _ := regexp.Compile(SAN)

	matched := r.FindStringSubmatch(moveToParse)

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
		//error finding the origin
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
		dest := mv.algebraic[2:4]
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
			if mv.algebraic[:2] == sq {
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
	return "", nil // TODO: this should really return an error
}

/*******************************************************************************

	Notation Stuff:

*******************************************************************************/

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
