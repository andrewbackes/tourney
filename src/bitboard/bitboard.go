/*

 Project: Tourney

 Module: bitboards
 Description: Support functions for working with bitbaords

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/14/2014

*/

package bitboard

import "fmt"

type BB struct {
	Value uint64
}

func (x BB) Print() {
	for i := 7; i >= 0; i-- {
		fmt.Printf("%08s\n", fmt.Sprintf("%b", (x.Value>>uint64(8*i)&255)))
	}
}

func Print(bb uint64) {
	var one uint64
	one = 1
	for i := 63; i >= 0; i-- {
		if bb&(one<<uint(i)) == uint64(0) {
			fmt.Print("0")
		} else {
			fmt.Print("1")
		}
		if i%8 == 0 {
			fmt.Print("\n")
		}
	}
}

/*
In C++
======

void bitprint(bitboard bb)
{
	bitboard x=1;
	for(int i=63; i>=0; i--){
		if( bb & (x<<i) )
			std::cout << "1";
		else
			std::cout << "0";
		if(i%8 == 0)
			std::cout << std::endl;
	}
	std::cout << std::endl;
}

*/
