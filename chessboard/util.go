package chessboard

func getUnicode(colour, piece int) string {
	switch piece {
	case Pawn:
		if colour == White {
			return "\u2659"
		}
		return "\u265F"
	case Rook:
		if colour == White {
			return "\u2656"
		}
		return "\u265C"
	case Knight:
		if colour == White {
			return "\u2658"
		}
		return "\u265E"
	case Bishop:
		if colour == White {
			return "\u2657"
		}
		return "\u265D"
	case Queen:
		if colour == White {
			return "\u2655"
		}
		return "\u265B"
	case King:
		if colour == White {
			return "\u2654"
		}
		return "\u265A"
	}
	return ""
}
