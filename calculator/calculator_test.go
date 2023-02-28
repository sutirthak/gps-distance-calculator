package calculator

import (
	"testing"
)
type speedtestCase struct {
	arg1 float64
	arg2 float64
	want float64
}

func TestCalculateAvarageSpeed(t *testing.T) {
	cases := []speedtestCase{
		{48,6,8},
		{32.3,2,16.15},
		{48,0,0},
	}

	for _, tc := range cases {
		got := CalculateAvarageSpeed(tc.arg1, tc.arg2)
		if tc.want != got {
			t.Errorf("Expected '%f', but got '%f'", tc.want, got)
		}
	}
}
