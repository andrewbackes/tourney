/*

 Project: Tourney

 Module: perft
 Description: debug for movegen

 TODO: turn this into a movegen package test

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*******************************************************************************

	Perft:

*******************************************************************************/

func perft(G Game, depth int) (uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	var nodes, checks, castles, mates, captures, promotions, enpassant uint64
	var moveCount uint64

	if depth == 0 {
		return 1, 0, 0, 0, 0, 0, 0
	}

	toMove := G.toMove()
	notToMove := []Color{BLACK, WHITE}[toMove]

	isChecked := G.isInCheck(toMove)
	ml := G.MoveGen()

	for _, mv := range ml {
		temp := G
		temp.MakeMove(mv)

		if temp.isInCheck(toMove) == false {
			//Count it for mate:
			moveCount += 1
			n, c, cstl, m, cap, p, enp := perft(temp, depth-1)
			nodes += n
			checks += c + toInt(temp.isInCheck(notToMove))
			castles += cstl + toInt(isCastle(&G, mv))
			mates += m
			captures += cap + toInt(isCapture(&G, mv))
			promotions += p + toInt(isPromotion(&G, mv))
			enpassant += enp + toInt(isEnPassant(&G, mv))
		}
	}
	if moveCount == 0 && isChecked {
		mates += 1
	}

	return nodes, checks, castles, mates, captures, promotions, enpassant

}

/*******************************************************************************

	Perft Suite:

*******************************************************************************/

func PerftSuite(filename string, maxdepth int) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
	}
	scanner := bufio.NewScanner(f)

	type Test struct {
		fen   string
		nodes []uint64
	}
	var test []Test
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")

		var newTest Test
		newTest.fen = words[0]
		newTest.nodes = append(newTest.nodes, 1) // depth 0 = 1 node

		for i := 1; i < len(words); i++ {
			n, _ := strconv.ParseUint(strings.Split(words[i], " ")[1], 10, 0)
			newTest.nodes = append(newTest.nodes, n)
		}

		test = append(test, newTest)
	}
	f.Close()

	for i, t := range test {
		var G Game
		G.LoadFEN(t.fen)
		fmt.Print("FEN ", i+1, ": ")
		for depth, nodes := range t.nodes {
			/*
				if nodes >= 100000000 {
					break
				}
			*/
			if depth > maxdepth {
				break
			}
			fmt.Print("D", depth, ": ")
			perftNodes, _, _, _, _, _, _ := perft(G, depth)
			fmt.Print(map[bool]string{true: "pass. ", false: "FAIL. "}[(perftNodes == nodes)])
		}
		fmt.Print("\n")
	}

}

/*******************************************************************************

	Divide:

*******************************************************************************/

func divide(G Game, depth int) {
	fmt.Println("Depth", depth)
	var nodes, moveCount uint64
	ml := G.MoveGen()
	toMove := G.toMove()
	for _, mv := range ml {
		temp := G
		temp.MakeMove(mv)

		if temp.isInCheck(toMove) == false {
			//Count it for mate:
			moveCount += 1
			n, _, _, _, _, _, _ := perft(temp, depth-1)
			fmt.Println(mv.algebraic, ":", n)
			nodes += n
		}
	}
	fmt.Println("Total: ", nodes, ". moves:", moveCount)
}

/*******************************************************************************

	Helpers:

*******************************************************************************/

func isCastle(G *Game, m Move) bool {
	from, _ := getIndex(m.algebraic)
	_, p := G.board.onSquare(from)
	if p == KING {
		if (m.algebraic == "e1g1") || (m.algebraic == "e1c1") || (m.algebraic == "e8g8") || (m.algebraic == "e8c8") {
			return true
		}
	}
	return false
}

func isCapture(G *Game, m Move) bool {
	_, to := getIndex(m.algebraic)
	_, cap := G.board.onSquare(to)
	return (cap != NONE)
}

func isPromotion(G *Game, m Move) bool {
	// TODO: will not work when more notation is added
	return (len(m.algebraic) > 4)
}

func isEnPassant(G *Game, m Move) bool {
	if G.enPassant == 64 {
		return false
	}
	from, to := getIndex(m.algebraic)
	_, p := G.board.onSquare(from)
	return (p == PAWN) && (to == G.enPassant) && ((from-to)%8 != 0)
}

func toInt(b bool) uint64 {
	if b == true {
		return 1
	} else {
		return 0
	}
}
