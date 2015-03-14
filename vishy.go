package vishycore

import (
	"errors"
)

type Board [12][12]int
type Color int

const (
	White Color = iota
	Black
)

type Pos struct {
	rank int
	file int
}

func NewBoard() Board {
	return [12][12]int{[12]int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99},
		[12]int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99},
		[12]int{99, 99, 1, 2, 3, 4, 5, 3, 2, 1, 99, 99},
		[12]int{99, 99, 6, 6, 6, 6, 6, 6, 6, 6, 99, 99},
		[12]int{99, 99, 0, 0, 0, 0, 0, 0, 0, 0, 99, 99},
		[12]int{99, 99, 0, 0, 0, 0, 0, 0, 0, 0, 99, 99},
		[12]int{99, 99, 0, 0, 0, 0, 0, 0, 0, 0, 99, 99},
		[12]int{99, 99, 0, 0, 0, 0, 0, 0, 0, 0, 99, 99},
		[12]int{99, 99, 7, 7, 7, 7, 7, 7, 7, 7, 99, 99},
		[12]int{99, 99, 8, 9, 10, 11, 12, 10, 9, 8, 99, 99},
		[12]int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99},
		[12]int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99},
	}
}

func (b Board) IsKingInCheck(kingColor Color) bool {
	var king, enemyQueen, enemyRook, enemyPawn, enemyBishop, enemyKnight int
	if kingColor == White {
		king = 4
		enemyQueen = 12
		enemyRook = 8
		enemyBishop = 10
		enemyPawn = 7
		enemyKnight = 9

	} else {
		king = 11
		enemyQueen = 5
		enemyRook = 1
		enemyPawn = 6
		enemyBishop = 3
		enemyKnight = 2
	}

	pos, err := findPiecePos(king, b)
	if err != nil {
		panic(err)
	}

	//check for pawn check
	if (kingColor == White && (b[pos.rank+1][pos.file-1] == enemyPawn || b[pos.rank+1][pos.file+1] == enemyPawn)) || (kingColor == Black && (b[pos.rank-1][pos.file-1] == enemyPawn || b[pos.rank-1][pos.file+1] == enemyPawn)) {
		return true
	}

	var i, piece int
	i = pos.rank + 1

	//go up the board
	for {
		piece = b[i][pos.file]
		if piece == 99 {
			break
		} else if piece == enemyQueen || piece == enemyRook {
			return true
		} else if piece != 0 {
			break
		}
		i++
	}

	//go down the board
	i = pos.rank - 1
	for {
		piece = b[i][pos.file]
		if piece == 99 {
			break
		} else if piece == enemyQueen || piece == enemyRook {
			return true
		} else if piece != 0 {
			break
		}
		i--
	}

	//go right
	i = pos.file - 1
	for {
		piece = b[pos.rank][i]
		if piece == 99 {
			break
		} else if piece == enemyQueen || piece == enemyRook {
			return true
		} else if piece != 0 {
			break
		}
		i--
	}

	//go left
	i = pos.file + 1
	for {
		piece = b[pos.rank][i]
		if piece == 99 {
			break
		} else if piece == enemyQueen || piece == enemyRook {
			return true
		} else if piece != 0 {
			break
		}
		i++
	}

	//top right
	i, j := pos.rank+1, pos.file-1
	for {
		piece = b[i][j]
		if piece == 99 {
			break
		} else if piece == enemyBishop || piece == enemyQueen {
			return true
		} else if piece != 0 {
			break
		}
		i++
		j--
	}

	//top left
	i, j = pos.rank+1, pos.file+1
	for {
		piece = b[i][j]
		if piece == 99 {
			break
		} else if piece == enemyBishop || piece == enemyQueen {
			return true
		} else if piece != 0 {
			break
		}
		i++
		j++
	}

	//bottom left
	i, j = pos.rank-1, pos.file+1
	for {
		piece = b[i][j]
		if piece == 99 {
			break
		} else if piece == enemyBishop || piece == enemyQueen {
			return true
		} else if piece != 0 {
			break
		}
		i--
		j++
	}

	//bottom right
	i, j = pos.rank-1, pos.file-1
	for {
		piece = b[i][j]
		if piece == 99 {
			break
		} else if piece == enemyBishop || piece == enemyQueen {
			return true
		} else if piece != 0 {
			break
		}
		i--
		j--
	}

	//move like a knight
	i, j = pos.rank, pos.file
	if b[i+2][j-1] == enemyKnight || b[i+2][j+1] == enemyKnight || b[i-2][j-1] == enemyKnight || b[i-2][j+1] == enemyKnight {
		return true
	}
	return false
}

func findPiecePos(piece int, b Board) (Pos, error) {
	for i := 2; i < 10; i++ {
		for j := 2; j < 10; j++ {
			if b[i][j] == piece {
				return Pos{rank: i, file: j}, nil
			}
		}
	}
	return Pos{}, errors.New("Piece not found")
}
