/*******************************************************************************

 Project: Tourney

 Module: notation
 Description: notation parsing utilities

 The internals of Tourney use Pure Coordinate Notation to represent moves. For
 example, a2a4 means to move what is on a2 to a4. However, UCI and WB require
 Standard Algebraic Notation. The purpose of this module is to translate between
 PCN and SAN.

 Author(s): Andrew Backes
 Created: 7/29/2014


 TODO:
 	-make PCNtoSAN()
 	-instead of calling LegalMoveGen(), should have a stripped down version
 	for just the type of piece moving.

*******************************************************************************/

package main

import (
	"errors"
	//"fmt"
	"regexp"
	"strings"
)

/*******************************************************************************

	PCN Convertion:

*******************************************************************************/

func ConvertToPCN(G *Game, moveToParse string) (string, error) {
	// Converts Standard Algebraic Notation (SAN) to Pure Coordinate Notation (PCN)
	// Examples of PCN are: e2e4 (and) e7e8Q

	// TODO: needs error handling.
	// TODO: what about promotion captures? or ambiguous promotions?
	//		 -Illegal move: f7g8 (raw: fxg8=Q).
	//		 -illegal move: move axb8=Q+

	// Check for null move:
	if moveToParse == "0000" {
		return moveToParse, nil
	}

	// Check for castling:
	if moveToParse == "O-O" {
		return []string{"e1g1", "e8g8"}[G.toMove()], nil
	}
	if moveToParse == "O-O-O" {
		return []string{"e1c1", "e8c8"}[G.toMove()], nil
	}

	// Strip uneeded characters:
	moveToParse = strings.Replace(moveToParse, "-", "", -1)

	// First check to see if it is already in the correct form.
	PCN := "([a-h][1-8])([a-h][1-8])([QBNRqbnr]?)"
	matches, _ := regexp.MatchString(PCN, moveToParse)
	if matches {
		parsed := moveToParse[:len(moveToParse)-1]
		// Some engines dont capitalize the promotion piece:
		parsed += strings.ToLower(moveToParse[len(moveToParse)-1:])
		// some engines dont specify the promotion piece, assume queen:
		if (parsed[1] == '7' && parsed[3] == '8') || (parsed[1] == '2' && parsed[3] == '1') {
			if len(parsed) <= 4 {
				f, _ := getIndex(parsed)
				_, p := G.Board.onSquare(f)
				if p == PAWN {
					parsed += "q"
				}
			}
		}
		return parsed, nil
	}

	//	    (piece)    (from)  (from)  (cap) (dest)      (promotion)        (chk  )
	SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x]?)([a-h][1-8])([=]?[BNPQRbnpqr]?)([+#]?)"
	r, _ := regexp.Compile(SAN)

	matched := r.FindStringSubmatch(moveToParse)
	if len(matched) == 0 {
		return moveToParse, errors.New("Error parsing move from engine: '" + moveToParse + "'")
	}

	piece := matched[1]
	fromFile := matched[2]
	fromRank := matched[3]
	//action := matched[4]      // capture or promote
	destination := matched[5] //or promotion piece if action="="
	//check := matched[6]       //or mate
	promote := strings.Replace(matched[6], "=", "", 1)

	if piece == "" {
		piece = "P"
	}

	origin, err := originOfPiece(piece, destination, fromFile, fromRank, G)
	if err != nil {
		//fmt.Println(err)
		//fmt.Println(G.FEN())
		//fmt.Println(moveToParse)
		//G.PrintHUD()
		return moveToParse, errors.New("Error finding source square of move: '" + moveToParse + "'.")
	}

	// Some engines dont tell you to promote to queen, so assume so in that case:
	/*if piece == "P" && ((origin[1] == '7' && destination[1] == '8') || (origin[1] == '2' && destination[1] == '1')) {
		if promote == "" {
			promote = "Q"
		}
	}
	*/
	return origin + destination + strings.ToLower(promote), nil
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
	bits := G.Board.PieceBB[color][pieceMap[piece]]
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
	//DEBUG:
	/*
		fmt.Println("params: ", piece, destination, fromFile, fromRank)
		fmt.Println("color: ", color)
		fmt.Println("legalMoves:", legalMoves)
		fmt.Println("eligableMoves:", eligableMoves)
		fmt.Println("eligableSquares:", eligableSquares)
	*/
	return "", errors.New("Notation: Can not find source square.")
}

/*******************************************************************************

	SAN Convertion:

*******************************************************************************/

func ConvertToSAN(G *Game, moveToParse string) (string, error) {

	return "", nil
}

/*******************************************************************************

	Notation Utilities:

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
		p = map[string]Piece{"Q": QUEEN, "N": KNIGHT, "B": BISHOP, "R": ROOK,
			"q": QUEEN, "n": KNIGHT, "b": BISHOP, "r": ROOK}
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

// Decides if a string is a chess move or not.
func isMove(s string) bool {

	if s == "O-O" || s == "O-O-O" || s == "0000" {
		return true
	}

	PCN := "^([a-h][1-8])([a-h][1-8])([QBNRqbnr]?)$"
	matches, _ := regexp.MatchString(PCN, s)
	if matches {
		return true
	}

	SAN := "^([BKNPQR]?)([a-h]?)([0-9]?)([x=]?)([BKNPQR]|[a-h][1-8])([+#]?)$"
	matches, _ = regexp.MatchString(SAN, s)
	if matches {
		return true
	}

	return false
}
