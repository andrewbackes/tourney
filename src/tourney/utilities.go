/*

 Project: Tourney

 Module: utilities
 Description: misc. functions and helpers and what not

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/28/2014

*/

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
)

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

/*******************************************************************************

	Data Verificaton:

*******************************************************************************/

func GetMD5(filepath string) (string, error) {
	// returns the MD5 sum of the file

	const filechunk = 8192

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)

		io.WriteString(hash, string(buf)) // append into the hash
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
