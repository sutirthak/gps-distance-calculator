package calculator

import (
	"fmt"
	"testing"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
)
type speedtestCase struct {
	totalDistance float64
	count float64
	want float64
}

type distancetestCase struct {
	start models.Coordinate
	end models.Coordinate
	want float64
}

func TestCalculateAvarageSpeed(t *testing.T) {
	cases := []speedtestCase{
		{48,6,8},
		{32.3,2,16.15},
		{0.00568,7,0.0008114285714285715},
		{12,0,0},
	}

	for _, tc := range cases {
		got := CalculateAvarageSpeed(tc.totalDistance, tc.count)
		fmt.Println(got)
		if tc.want != got {
			t.Errorf("Expected '%f', but got '%f'", tc.want, got)
		}
	}
}

func TestCalculateDistanceInMeter(t *testing.T) {
	cases := []distancetestCase{
		{models.Coordinate{Latitude:1.360495,Longitude:103.954334},models.Coordinate{ Latitude:1.320491,Longitude:103.964334},4585.041981327221},
		{models.Coordinate{Latitude:51.512722,Longitude:-0.288552},models.Coordinate{ Latitude:51.516100,Longitude:0.068025},24677.456240403306},
		{models.Coordinate{Latitude:53.478612,Longitude:6.250578},models.Coordinate{ Latitude:50.752342,Longitude:5.916981},304001.0210460888},

	}
	for _, tc := range cases {
		got := CalculateDistanceInMeter(tc.start, tc.end)
		if tc.want != got {
			t.Errorf("Expected '%f', but got '%f'", tc.want, got)
		}
	}
}

func TestDistanceIsTheSameInBothDirections(t *testing.T){
	startPosition:=models.Coordinate{Latitude: -33.926510,Longitude:18.364603}
	endPosition:=models.Coordinate{Latitude: -26.208450,Longitude:28.040572}
	distance1:=CalculateDistanceInMeter(startPosition,endPosition)
	distance2:=CalculateDistanceInMeter(endPosition,startPosition)
	if distance1!=distance2{
		t.Errorf("distance should be equal")
}
}
