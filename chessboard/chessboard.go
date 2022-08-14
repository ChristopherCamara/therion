package chessboard

import "fmt"

// Bitboard is an alias for uint64
type Bitboard uint64

type square int

//The 8 different bitboards contained in a chessboard
const (
	White = iota
	Black
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

//the 64 squares of the chessboard - in Little Endian Rank-File Mapping order
const (
	A1 square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// Chessboard contains the various bitboards
type Chessboard struct {
	bitboards [8]Bitboard
}

// New chessboard initialzed to starting position
func New() Chessboard {
	var c Chessboard
	c.bitboards = [8]Bitboard{
		0x000000000000FFFF,
		0xFFFF000000000000,
		0x00FF00000000FF00,
		0x8100000000000081,
		0x4200000000000042,
		0x2400000000000024,
		0x1000000000000010,
		0x0800000000000008,
	}
	return c
}

//Print out the Chessboard with the unicode of the pieces
func (c Chessboard) Print() {
	for i := H8; i >= A1; i-- {
		printed := false
		for j := 0; j < 2; j++ {
			for k := 2; k < 8; k++ {
				if c.bitboards[j]&c.bitboards[k]>>i&1 == 1 {
					fmt.Printf("%s ", getUnicode(j, k))
					printed = true
				}
			}
		}
		if !printed {
			fmt.Print(". ")
		}
		if i%8 == 0 {
			fmt.Println()
		}
	}
}
