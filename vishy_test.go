package vishycore

import "testing"

func TestFENToBoardStateValidCases(t *testing.T) {
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
	}
	for _, c := range cases {
		got, err := NewBoardStateFromFen(c.in)
		if err != nil {
			t.Errorf("in: %s\n want: %+v\n error: %s\n", c.in, c.want, err)
		}
		if got != c.want {
			t.Errorf("in: %s\n want: %+v\n, got: %+v\n", c.in, c.want, got)
		}
	}
}
