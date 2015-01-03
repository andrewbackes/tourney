package main

import (
	"testing"
)

func TestConvertToPCN(t *testing.T) {

	tests := [][]string{
		{
			// 0
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"e2e4",
			"e2e4",
		},
		{
			// 1
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"e4",
			"e2e4",
		},
		{
			// 2
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"b1c3",
			"b1c3",
		},
		{
			// 3
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"Nc3",
			"b1c3",
		},
		{
			// 4
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"c3",
			"c2c3",
		},
		{
			// 5
			"2r5/4R2P/p7/8/3k4/1P4P1/P4PK1/8 w - - 0 47",
			"h8Q",
			"h7h8q",
		},
		{
			// 6
			"2r5/4R2P/p7/8/3k4/1P4P1/P4PK1/8 w - - 0 47",
			"h7h8q",
			"h7h8q",
		},
		{
			// 7
			"r1bqk2r/pppn1pbp/3p2p1/5p2/2PP4/2N1P1P1/PP3PBP/R2QK1NR w KQkq - 1 8", // FEN
			"Nge2", // SAN
			"g1e2", // PCN
		},
		{
			// 8
			"1q4k1/P5p1/1K6/2pb4/7B/8/8/8 w - - 0 145",
			"axb8=Q+",
			"a7b8q",
		},
		{
			// 9
			"3r3k/4P3/7K/8/8/3p2P1/7P/8 w - - 0 69",
			"exd8=Q#",
			"e7d8q",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"Qh3",
			"f3h3",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"Nb5",
			"c3b5",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"Pxh3",
			"g2h3",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"xh3",
			"g2h3",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"O-O-O",
			"e1c1",
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			"O-O",
			"e1g1",
		},
		{
			"8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1",
			"gf4",
			"g3f4",
		},
		{
			"8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1",
			"ef4",
			"e3f4",
		},
		{
			"8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1",
			"gxf4",
			"g3f4",
		},
		{
			// 19
			"8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1",
			"Rf4",
			"b4f4",
		},
		{
			// 20
			"8/2p5/3p4/KP5k/1R3p1r/4P1P1/8/8 w - - 0 1",
			"Rxf4",
			"b4f4",
		},
	}

	for i, _ := range tests {
		tester := NewGame()
		tester.LoadFEN(tests[i][0])
		SAN := tests[i][1]
		PCN := tests[i][2]
		answer, err := ConvertToPCN(&tester, SAN)
		if answer != PCN {
			t.Error(i, "FAILED: Expected", PCN, "got", answer, "-", err)
		} else {
			//t.Log(i, "PASSED. ")
		}

	}
}
