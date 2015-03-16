package vishycore

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Board [12][12]int
type Color int

type BoardState struct {
	board                 Board
	sideToMove            rune
	castlingAbility       string
	enPassantTargetSquare string
	halfMoveClock         int
	fullMoveCounter       int
}

func NewBoardState() BoardState {
	return BoardState{
		board:                 NewBoard(),
		sideToMove:            'w',
		castlingAbility:       "KQkq",
		enPassantTargetSquare: "-",
		halfMoveClock:         0,
		fullMoveCounter:       1,
	}
}

func NewBoardStateFromFen(fen string) (BoardState, error) {
	var boardState BoardState
	throwError := func() (BoardState, error) {
		return NewBoardState(), errors.New("Invalid FEN")
	}
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
		return throwError()
	}

	/*
		<Piece Placement> ::= <rank8>'/'<rank7>'/'<rank6>'/'<rank5>'/'<rank4>'/'<rank3>'/'<rank2>'/'<rank1>
	*/
	piecePlacement := components[0]
	ranks := strings.Split(piecePlacement, "/")
	if len(ranks) != 8 {
		return throwError()
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
				return throwError()
			}
			switch runeVal {
			case '8':
				if len(rank) != 1 {
					return throwError()
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
				return throwError()
			}
		}
	}

	//place pieces on the board
	for i, rank := range ranks {
		j := 0
		for _, runeVal := range rank {
			//Now I realize the power of closures ;)
			fillEmptySpaces := func(x int) {
				m := x
				for m > 0 {
					board[9-i][9-j] = 0
					m--
					j++
				}
			}
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
		return throwError()
	} else if side := sideToMove[0]; side != 'w' && side != 'b' {
		return throwError()
	} else if isBlackKingInCheck := IsKingInCheck(Black, board); side == 'w' && isBlackKingInCheck {
		return throwError()
	} else if isWhiteKingInCheck := IsKingInCheck(White, board); side == 'b' && isWhiteKingInCheck {
		return throwError()
	} else if isWhiteKingInCheck && isBlackKingInCheck {
		return throwError()
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
					return throwError()
				}
			case 'K':
				if encounter['K'] || board[2][5] != pieceMap['K'] || board[2][2] != pieceMap['R'] {
					return throwError()
				}
				encounter['K'] = true
			case 'Q':
				if encounter['Q'] || board[2][5] != pieceMap['K'] || board[2][9] != pieceMap['R'] {
					return throwError()
				}
				encounter['Q'] = true
			case 'k':
				if encounter['k'] || board[9][5] != pieceMap['k'] || board[9][2] != pieceMap['r'] {
					return throwError()
				}
				encounter['k'] = true
			case 'q':
				if encounter['q'] || board[9][5] != pieceMap['k'] || board[9][9] != pieceMap['r'] {
					return throwError()
				}
				encounter['q'] = true
			default:
				return throwError()
			}
		}
	} else {
		return throwError()
	}

	/*
		<En passant target square> ::= '-' | <epsquare>
		<epsquare>   ::= <fileLetter> <eprank>
		<fileLetter> ::= 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h'
		<eprank>     ::= '3' | '6'
	*/

	enPassantTargetSquare := components[3]
	eLength := len(enPassantTargetSquare)

	if eLength == 1 {
		if enPassantTargetSquare[0] != '-' {
			return throwError()
		}
	} else if eLength == 2 {
		f1 := false
		validate := func(algNotation string, pawnType rune) bool {
			var expectedSide rune
			if unicode.IsUpper(pawnType) {
				expectedSide = 'b'
			} else {
				expectedSide = 'w'
			}

			pos, err := getPos(algNotation)
			if err != nil || expectedSide != rune(sideToMove[0]) || board[pos.x][pos.y] != pieceMap[pawnType] || board[pos.x-1][pos.y] != 0 || board[pos.x-2][pos.y] != 0 {
				return true
			}
			return false
		}
		switch enPassantTargetSquare {
		case "a3":
			f1 = validate("a4", 'P')
		case "b3":
			f1 = validate("b4", 'P')
		case "c3":
			f1 = validate("c4", 'P')
		case "d3":
			f1 = validate("d4", 'P')
		case "e3":
			f1 = validate("e4", 'P')
		case "f3":
			f1 = validate("f4", 'P')
		case "g3":
			f1 = validate("g4", 'P')
		case "a6":
			f1 = validate("a5", 'p')
		case "b6":
			f1 = validate("b5", 'p')
		case "c6":
			f1 = validate("c5", 'p')
		case "d6":
			f1 = validate("d5", 'p')
		case "e6":
			f1 = validate("e5", 'p')
		case "f6":
			f1 = validate("f5", 'p')
		case "g6":
			f1 = validate("g5", 'p')
		case "h6":
			f1 = validate("h5", 'p')
		default:
			f1 = true
		}
		if f1 {
			return throwError()
		}
	} else {
		return throwError()
	}

	/*
		<Halfmove Clock> ::= <digit> {<digit>}
		<digit> ::= '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
	*/
	halfMoveClock := components[4]
	if len(halfMoveClock) < 1 {
		return throwError()
	} else {
		for _, val := range halfMoveClock {
			if !unicode.IsDigit(val) {
				return throwError()
			}
		}
	}

	/*
		<Fullmove counter> ::= <digit19> {<digit>}
		<digit19> ::= '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
		<digit>   ::= '0' | <digit19>
	*/

	fullMoveCounter := components[5]
	fLength := len(fullMoveCounter)

	if fLength < 1 {
		return throwError()
	}
	switch fullMoveCounter[0] {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		return throwError()
	}
	for _, val := range fullMoveCounter[1:] {
		switch val {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return throwError()
		}
	}
	halfMoveClockInt, _ := strconv.Atoi(halfMoveClock)
	fullMoveCounterInt, _ := strconv.Atoi(fullMoveCounter)

	boardState = BoardState{
		board:                 board,
		sideToMove:            rune(sideToMove[0]),
		castlingAbility:       castlingAbility,
		enPassantTargetSquare: enPassantTargetSquare,
		halfMoveClock:         halfMoveClockInt,
		fullMoveCounter:       fullMoveCounterInt,
	}
	return boardState, nil
}

const (
	White Color = iota
	Black
)

type Pos struct {
	x int
	y int
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

func getPos(algebNotation string) (Pos, error) {
	file, rank := algebNotation[0], algebNotation[1]

	pos := Pos{}

	switch file {
	case 'a':
		pos.y = 9
	case 'b':
		pos.y = 8
	case 'c':
		pos.y = 7
	case 'd':
		pos.y = 6
	case 'e':
		pos.y = 5
	case 'f':
		pos.y = 4
	case 'g':
		pos.y = 3
	case 'h':
		pos.y = 2
	default:
		return pos, errors.New("Invalid algebraic notation")

	}

	switch rank {
	case '1':
		pos.x = 2
	case '2':
		pos.x = 3
	case '3':
		pos.x = 4
	case '4':
		pos.x = 5
	case '5':
		pos.x = 6
	case '6':
		pos.x = 7
	case '7':
		pos.x = 8
	case '8':
		pos.x = 9
	default:
		return pos, errors.New("Invalid algebraic notation")
	}

	return pos, nil
}
func IsKingInCheck(kingColor Color, b Board) bool {
	var king, enemyQueen, enemyRook, enemyPawn, enemyBishop, enemyKnight int
	if kingColor == White {
		king = pieceMap['K']
		enemyQueen = pieceMap['q']
		enemyRook = pieceMap['r']
		enemyBishop = pieceMap['b']
		enemyPawn = pieceMap['p']
		enemyKnight = pieceMap['n']

	} else {
		king = pieceMap['k']
		enemyQueen = pieceMap['Q']
		enemyRook = pieceMap['R']
		enemyPawn = pieceMap['P']
		enemyBishop = pieceMap['B']
		enemyKnight = pieceMap['N']
	}

	pos, err := findPiecePos(king, b)
	if err != nil {
		panic(err)
	}

	//check for pawn check
	if (kingColor == White && (b[pos.x+1][pos.y-1] == enemyPawn || b[pos.x+1][pos.y+1] == enemyPawn)) || (kingColor == Black && (b[pos.x-1][pos.y-1] == enemyPawn || b[pos.x-1][pos.y+1] == enemyPawn)) {
		return true
	}

	var i, piece int
	i = pos.x + 1

	//go up the board
	for {
		piece = b[i][pos.y]
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
	i = pos.x - 1
	for {
		piece = b[i][pos.y]
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
	i = pos.y - 1
	for {
		piece = b[pos.x][i]
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
	i = pos.y + 1
	for {
		piece = b[pos.x][i]
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
	i, j := pos.x+1, pos.y-1
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
	i, j = pos.x+1, pos.y+1
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
	i, j = pos.x-1, pos.y+1
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
	i, j = pos.x-1, pos.y-1
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
	i, j = pos.x, pos.y
	if b[i+2][j-1] == enemyKnight || b[i+2][j+1] == enemyKnight || b[i-2][j-1] == enemyKnight || b[i-2][j+1] == enemyKnight {
		return true
	}
	return false
}

func findPiecePos(piece int, b Board) (Pos, error) {
	for i := 2; i < 10; i++ {
		for j := 2; j < 10; j++ {
			if b[i][j] == piece {
				return Pos{x: i, y: j}, nil
			}
		}
	}
	return Pos{}, errors.New("Piece not found")
}
