package vishycore

import (
	"unicode"
)

type Square string
type Piece rune

type Turn int

const (
	WhiteTurn Turn = iota
	BlackTurn
)

type Color int

const (
	White Color = iota
	Black
)

type BoardState struct {
	Board                       [8][8]Piece
	Turn                        Turn
	EnPassantTargetSquare       Square
	IsBlackQSideCastleAvailable bool
	IsWhiteQSideCastleAvailable bool
	IsBlackKSideCastleAvailable bool
	IsWhiteKSideCastleAvailable bool
	HalfMoves                   int
	FullMoves                   int
}

func CreateNewBoard() BoardState {
	boardState := BoardState{
		IsBlackQSideCastleAvailable: true,
		IsWhiteQSideCastleAvailable: true,
		IsBlackKSideCastleAvailable: true,
		IsWhiteKSideCastleAvailable: true,
		EnPassantTargetSquare:       "-",
		HalfMoves:                   0,
		FullMoves:                   1,
		Turn:                        WhiteTurn,
		Board: [8][8]Piece{
			[8]Piece{'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r'},
			[8]Piece{'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p'},
			[8]Piece{},
			[8]Piece{},
			[8]Piece{},
			[8]Piece{},
			[8]Piece{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'},
			[8]Piece{'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R'},
		},
	}
	return boardState
}

func IsKingInCheck(boardState BoardState, kingColor Color) bool {
	var king Piece
	switch kingColor {
	case White:
		king = 'k'
	case Black:
		king = 'K'
	}

	//Searching the king, need a better implementation instead of searching for him
	var i, j int
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			if boardState.Board[i][j] == king {
				break
			}
		}
	}

	checkPath := func(k int, l int, pieceType []Piece) bool {
		piece := boardState.Board[k][l]
		if piece != 0 {
			var pieceColor Color
			if unicode.IsLower(rune(piece)) {
				pieceColor = Black
			} else {
				pieceColor = White
			}

			if pieceColor != kingColor {
				piecesLength := len(pieceType)
				for i := 0; i < piecesLength; i++ {
					if piece == pieceType[i] {
						return true
					}
				}
			}
		}
		return false
	}

	//look for attackers up the board
	for k := i + 1; k < 8; k++ {
		if checkPath(k, j, []Piece{'q', 'r'}) {
			return true
		}
	}

	//look for attackers down the board
	for l := i - 1; l >= 0; l-- {
		if checkPath(l, j, []Piece{'q', 'r'}) {
			return true
		}
	}

	return false
}
