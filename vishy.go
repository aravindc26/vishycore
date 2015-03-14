package vishycore

import (
	"errors"
	"strings"
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

var pieceMap = map[rune]int{
	'R': 1,
	'N': 2,
	'B': 3,
	'K': 4,
	'Q': 5,
	'P': 6,
	'p': 7,
	'r': 8,
	'n': 9,
	'b': 10,
	'k': 11,
	'q': 12,
}

func NewBoard() Board {
	return [12][12]int{
		[12]int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99},
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

func CreateBoardFromFen(fen string) (Board, error) {
	board := NewBoard()
	/*
			<FEN> ::=  <Piece Placement>
		       ' ' <Side to move>
		       ' ' <Castling ability>
		       ' ' <En passant target square>
		       ' ' <Halfmove clock>
		       ' ' <Fullmove counter>
	*/
	fen = strings.Trim(fen, " ")
	components := strings.Split(fen, " ")

	if len(components) != 6 {
		return board, errors.New("Invalid FEN")
	}

	/*
		<Piece Placement> ::= <rank8>'/'<rank7>'/'<rank6>'/'<rank5>'/'<rank4>'/'<rank3>'/'<rank2>'/'<rank1>
	*/
	piecePlacement := components[0]
	ranks := strings.Split(piecePlacement, "/")
	if len(ranks) != 8 {
		return board, errors.New("Invalid FEN")
	}
	/*
		<ranki>       ::= [<digit17>]<piece> {[<digit17>]<piece>} [<digit17>] | '8'
		<piece>       ::= <white Piece> | <black Piece>
		<digit17>     ::= '1' | '2' | '3' | '4' | '5' | '6' | '7'
		<white Piece> ::= 'P' | 'N' | 'B' | 'R' | 'Q' | 'K'
		<black Piece> ::= 'p' | 'n' | 'b' | 'r' | 'q' | 'k'
	*/
	for _, rank := range ranks {
		var sum int
		for _, runeVal := range rank {
			if sum > 8 {
				return board, errors.New("Invalid FEN")
			}
			switch runeVal {
			case '8':
				if len(rank) != 1 {
					return board, errors.New("Invalid FEN")
				}
			case '7':
				sum += 7
			case '6':
				sum += 6
			case '5':
				sum += 5
			case '4':
				sum += 4
			case '3':
				sum += 3
			case '2':
				sum += 2
			case '1', 'r', 'n', 'b', 'k', 'q', 'p', 'R', 'N', 'B', 'K', 'Q', 'P':
				sum += 1
			default:
				return board, errors.New("Invalid FEN")
			}
		}
	}

	//place pieces on the board
	for i, rank := range ranks {
		j := 0
		//Now I realize the power of closures ;)
		fillEmptySpaces := func(x int) {
			for ; j < x; j++ {
				board[9-i][9-j] = 0
			}
		}
		for _, runeVal := range rank {
			switch runeVal {
			case '8':
				fillEmptySpaces(8)
			case '7':
				fillEmptySpaces(7)
			case '6':
				fillEmptySpaces(6)
			case '5':
				fillEmptySpaces(5)
			case '4':
				fillEmptySpaces(4)
			case '3':
				fillEmptySpaces(3)
			case '2':
				fillEmptySpaces(2)
			case '1':
				fillEmptySpaces(1)
			default:
				board[9-i][9-j] = pieceMap[runeVal]
				j++
			}
		}
	}

	/*
		<Side to move> ::= {'w' | 'b'}
	*/

	sideToMove := components[1]
	if len(sideToMove) != 1 {
		return board, errors.New("Invalid FEN")
	} else if side := sideToMove[0]; side != 'w' || side != 'b' {
		return board, errors.New("Invalid FEN")
	} else if isBlackKingInCheck := IsKingInCheck(Black, board); side == 'w' && isBlackKingInCheck {
		return board, errors.New("Invalid FEN")
	} else if isWhiteKingInCheck := IsKingInCheck(White, board); side == 'b' && isWhiteKingInCheck {
		return board, errors.New("Invalid FEN")
	} else if isWhiteKingInCheck && isBlackKingInCheck {
		return board, errors.New("Invalid FEN")
	}

	/*
		<Castling ability> ::= '-' | ['K'] ['Q'] ['k'] ['q'] (1..4)
	*/
	castlingAbility := components[2]

	if cLength := len(castlingAbility); cLength >= 1 && cLength <= 4 {
		encounter := map[rune]bool{
			'k': false,
			'K': false,
			'q': false,
			'Q': false,
		}
		for _, val := range castlingAbility {
			switch val {
			case '_':
				if cLength != 1 {
					return board, errors.New("Invalid FEN")
				}
			case 'K':
				if encounter['K'] || board[2][5] != pieceMap['K'] || board[2][2] != pieceMap['R'] {
					return board, errors.New("Invalid FEN")
				}
				encounter['K'] = true
			case 'Q':
				if board[2][5] != pieceMap['K'] || board[2][9] != pieceMap['R'] {
					return board, errors.New("Invalid FEN")
				}
				encounter['Q'] = true
			case 'k':
				if board[9][5] != pieceMap['k'] || board[9][2] != pieceMap['r'] {
					return board, errors.New("Invalid FEN")
				}
				encounter['k'] = true
			case 'q':
				if board[9][5] != pieceMap['k'] || board[9][9] != pieceMap['r'] {
					return board, errors.New("Invalid FEN")
				}
				encounter['q'] = true
			default:
				return board, errors.New("Invalid FEN")
			}
		}
	} else {
		return board, errors.New("Invalid FEN")
	}
	return board, nil
}

func IsKingInCheck(kingColor Color, b Board) bool {
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
