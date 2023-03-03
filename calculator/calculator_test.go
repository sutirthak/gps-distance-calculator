package calculator

import (
	"testing"

	"github.com/mursalinsk-qi/gps-distance-calculator/models"
	"github.com/stretchr/testify/assert"
)

type speedtestCase struct {
	tripData      models.Trip
	current_speed float64
	expected      float64
}

type distancetestCase struct {
	start    models.Position
	end      models.Position
	expected float64
}

func TestCalculateAvarageSpeed(t *testing.T) {
	cases := []speedtestCase{
		{models.Trip{AvgSpeed: 50, DataCount: 2}, 29, 43},
		{models.Trip{AvgSpeed: 0.57, DataCount: 2}, 13.8, 4.98},
	}

	for _, tc := range cases {
		got := CalculateAvarageSpeed(tc.tripData, tc.current_speed)
		assert.Equal(t, tc.expected, got, "test case failed for avarage speed calculation")
	}
}

func TestCalculateDistanceInMeter(t *testing.T) {
	cases := []distancetestCase{
		{models.Position{Latitude: 1.360495, Longitude: 103.954334}, models.Position{Latitude: 1.320491, Longitude: 103.964334}, 4585.041981327221},
		{models.Position{Latitude: 51.512722, Longitude: -0.288552}, models.Position{Latitude: 51.516100, Longitude: 0.068025}, 24677.456240403306},
		{models.Position{Latitude: 53.478612, Longitude: 6.250578}, models.Position{Latitude: 50.752342, Longitude: 5.916981}, 304001.0210460888},
	}
	for _, tc := range cases {
		got := CalculateDistanceInMeter(tc.start, tc.end)
		assert.Equal(t, tc.expected, got, "test case failed for distance calculation")
	}
}

func TestDistanceIsTheSameInBothDirections(t *testing.T) {
	startPosition := models.Position{Latitude: -33.926510, Longitude: 18.364603}
	endPosition := models.Position{Latitude: -26.208450, Longitude: 28.040572}
	distance1 := CalculateDistanceInMeter(startPosition, endPosition)
	distance2 := CalculateDistanceInMeter(endPosition, startPosition)
	assert.Equal(t,distance1,distance2, "distance in both direction should be same")
}
