/*

 Project: Tourney

 Module: utilities
 Description: misc. functions and helpers and what not

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/28/2014

*/

package main

import "fmt"

/*******************************************************************************

	Bit Stuff:

*******************************************************************************/

// TODO: These are horribly inefficient functions. Help a brotha' out!

func popcount(b uint64) uint {
	var count uint
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			count += 1
		}
	}
	return count
}

func bitscan(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func BSF(b uint64) uint {
	for i := uint(0); i < 64; i++ {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	return 64
}

func BSR(b uint64) uint {
	for i := uint(63); i > 0; i-- {
		if (b & (1 << i)) != 0 {
			return i
		}
	}
	if b&1 != 0 {
		return 0
	}
	return 64
}

func bitprint(x uint64) {
	for i := 7; i >= 0; i-- {
		fmt.Printf("%08b\n", (x >> uint64(8*i) & 255))
	}
}
