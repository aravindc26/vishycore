package vishycore

import "testing"

func TestFENToBoardStateValidCases(t *testing.T) {
	board2 := NewBoard()
	board2[5][5] = pieceMap['P']
	board2[3][5] = 0
	board3 := NewBoard()
	board3[5][5] = pieceMap['P']
	board3[3][5] = 0
	board3[6][7] = pieceMap['p']
	board3[8][7] = 0
	board4 := NewBoard()
	board4[5][5] = pieceMap['P']
	board4[3][5] = 0
	board4[6][7] = pieceMap['p']
	board4[8][7] = 0
	board4[4][4] = pieceMap['N']
	board4[2][3] = 0
	cases := []struct {
		in   string
		want BoardState
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			BoardState{
				board:                 NewBoard(),
				sideToMove:            'w',
				castlingAbility:       "KQkq",
				enPassantTargetSquare: "-",
				halfMoveClock:         0,
				fullMoveCounter:       1,
			},
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			BoardState{
				board:                 board2,
				sideToMove:            'b',
				castlingAbility:       "KQkq",
				enPassantTargetSquare: "e3",
				halfMoveClock:         0,
				fullMoveCounter:       1,
			},
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			BoardState{
				board:                 board3,
				sideToMove:            'w',
				castlingAbility:       "KQkq",
				enPassantTargetSquare: "c6",
				halfMoveClock:         0,
				fullMoveCounter:       2,
			},
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			BoardState{
				board:                 board4,
				sideToMove:            'b',
				castlingAbility:       "KQkq",
				enPassantTargetSquare: "-",
				halfMoveClock:         1,
				fullMoveCounter:       2,
			},
		},
	}
	for _, c := range cases {
		got, err := NewBoardStateFromFen(c.in)
		if err != nil {
			t.Errorf("in: %s\n error: %s\n", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("in: %s\n want: %+v\n, got: %+v\n", c.in, c.want, got)
		}
	}
}
