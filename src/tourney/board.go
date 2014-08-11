/*

 Project: Tourney

 Module: boards
 Description: holds the board object and methods for interacting with it.
 	Bitboards, 'mailbox' view of the board, etc

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

import (
	"fmt"
)

//import "bitboard"

type Square uint8
type Piece uint8

const (
	PAWN Piece = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	NONE
)

// It sucks that this takes up so much room after formatted correctly:
const (
	H1 = iota
	G1
	F1
	E1
	D1
	C1
	B1
	A1
	H2
	G2
	F2
	E2
	D2
	C2
	B2
	A2
	H3
	G3
	F3
	E3
	D3
	C3
	B3
	A3
	H4
	G4
	F4
	E4
	D4
	C4
	B4
	A4
	H5
	G5
	F5
	E5
	D5
	C5
	B5
	A5
	H6
	G6
	F6
	E6
	D6
	C6
	B6
	A6
	H7
	G7
	F7
	E7
	D7
	C7
	B7
	A7
	H8
	G8
	F8
	E8
	D8
	C8
	B8
	A8
)

type Board struct {
	// contains all of the board state information
	// this includes all bitboards

	// it was turning into too much of a pain to use the type BB:
	pieceBB [2][6]uint64 //[player][piece]

}

func (B *Board) Print() {
	//TODO: maybe this should be a method of Game since we will want to also print
	//		game specific data. Such as moves to go, timers, etc.

	abbrev := [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}
	fmt.Println("+---+---+---+---+---+---+---+---+")
	for i := 1; i <= 64; i++ {
		square := uint(64 - i)

		fmt.Print("|")
		blankSquare := true
		for j := PAWN; j <= KING; j = j + 1 {
			for color := Color(0); color <= BLACK; color++ {
				if ((1 << square) & B.pieceBB[color][j]) != 0 {
					fmt.Print(" ", abbrev[color][j], " ")
					blankSquare = false
				}
			}
		}
		if blankSquare == true {
			fmt.Print("   ")
		}
		if square%8 == 0 {
			fmt.Println("|")
			fmt.Println("+---+---+---+---+---+---+---+---+")
		}
	}

}

func (B *Board) Clear() {
	B.pieceBB = [2][6]uint64{}
}

func (B *Board) Reset() {
	// puts the pieces in their starting/newgame positions
	for color := uint(0); color < 2; color = color + 1 {
		//Pawns first:
		B.pieceBB[color][PAWN] = 255 << (8 + (color * 8 * 5))
		//Then the rest of the pieces:
		B.pieceBB[color][KNIGHT] = (1 << (B1 + (color * 8 * 7))) ^ (1 << (G1 + (color * 8 * 7)))
		B.pieceBB[color][BISHOP] = (1 << (C1 + (color * 8 * 7))) ^ (1 << (F1 + (color * 8 * 7)))
		B.pieceBB[color][ROOK] = (1 << (A1 + (color * 8 * 7))) ^ (1 << (H1 + (color * 8 * 7)))
		B.pieceBB[color][QUEEN] = (1 << (D1 + (color * 8 * 7)))
		B.pieceBB[color][KING] = (1 << (E1 + (color * 8 * 7)))
	}
}

// Helpers:

func (B *Board) onSquare(square uint8) (Color, Piece) {
	// returns the piece on that square
	for c := WHITE; c <= BLACK; c++ {
		for p := PAWN; p <= KING; p++ {
			if (B.pieceBB[c][p] & (1 << square)) != 0 {
				return c, p
			}
		}
	}
	return NEITHER, NONE
}

func (B *Board) Occupied(c Color) uint64 {
	var mask uint64
	for p := PAWN; p <= KING; p++ {
		if c == BOTH {
			mask |= B.pieceBB[WHITE][p] | B.pieceBB[BLACK][p]
		} else {
			mask |= B.pieceBB[c][p]
		}
	}
	return mask
}
