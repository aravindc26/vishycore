package vishycore

import "testing"

func TestFENToBoardStateValidCases(t *testing.T) {
	board2 := NewBoard()
	board2[5][5] = pieceMap['p']
	board2[3][5] = 0
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
	}
	for _, c := range cases {
		got, err := NewBoardStateFromFen(c.in)
		if err != nil {
			t.Errorf("in: %s\n error: %s\n", c.in, err)
		}
		if got != c.want {
			t.Errorf("in: %s\n want: %+v\n, got: %+v\n", c.in, c.want, got)
		}
	}
}
