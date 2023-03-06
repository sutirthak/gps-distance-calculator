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

type validPositiontestCase struct {
	position    models.Position
	expected bool
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

// Testing for valid positions
func TestCalculateDistanceValidPosition(t *testing.T) {
	cases := []distancetestCase{
		{models.Position{Latitude: 1.360495, Longitude: 103.954334,IsValid: true}, models.Position{Latitude: 1.320491, Longitude: 103.964334,IsValid: true}, 4585.041981327221},
		{models.Position{Latitude: 51.512722, Longitude: -0.288552,IsValid: true}, models.Position{Latitude: 51.516100, Longitude: 0.068025,IsValid: true}, 24677.456240403306},
		{models.Position{Latitude: 53.478612, Longitude: 6.250578,IsValid: true}, models.Position{Latitude: 50.752342, Longitude: 5.916981,IsValid: true}, 304001.0210460888},
	}
	for _, tc := range cases {
		result,err := CalculateDistanceInMeter(tc.start, tc.end)
		assert.Nil(t, err,"error must be nil for valid position")
		assert.Equal(t, tc.expected, result, "test case failed for distance calculation")
	}
}
// Testing for invalid positions
func TestCalculateDistanceInvalidPosition(t *testing.T) {
	cases := []distancetestCase{
		{models.Position{Latitude: 0, Longitude: 0,IsValid: false}, models.Position{Latitude: 1.320491, Longitude: 103.964334,IsValid: true},-1},
		{models.Position{Latitude: 84.56321, Longitude: -183.58641,IsValid: false}, models.Position{Latitude: 24.45321, Longitude: 103.423561,IsValid: true},-1},
		{models.Position{Latitude: -98.56321, Longitude: 156.123456,IsValid: false}, models.Position{Latitude: 13.320491, Longitude: 109.964334,IsValid: true},-1},
		{models.Position{Latitude: -93.56324, Longitude: -183.2354,IsValid: false}, models.Position{Latitude: 1.320491, Longitude: 103.964334,IsValid: true},-1},
		{models.Position{Latitude: 64.32142, Longitude: -110.98632,IsValid: true}, models.Position{Latitude: 45, Longitude: 184,IsValid: false},-1},
		{models.Position{Latitude: -10.32142, Longitude: 108.98632,IsValid: true}, models.Position{Latitude: 99.320491, Longitude: -120.964334,IsValid: false},-1},
		{models.Position{Latitude: 64.32142, Longitude: 108.98632,IsValid: true}, models.Position{Latitude: -92.320491, Longitude: 183.964334,IsValid: false},-1},

	}
	for _, tc := range cases {
		result,err := CalculateDistanceInMeter(tc.start, tc.end)
		assert.NotNil(t,err,"there must be an error for invalid position")
		assert.Equal(t, tc.expected, result, "test case failed for distance calculation")
	}
}

func TestDistanceIsTheSameInBothDirections(t *testing.T) {
	startPosition := models.Position{Latitude: -33.926510, Longitude: 18.364603}
	endPosition := models.Position{Latitude: -26.208450, Longitude: 28.040572}
	distance1,_ := CalculateDistanceInMeter(startPosition, endPosition)
	distance2,_ := CalculateDistanceInMeter(endPosition, startPosition)
	assert.Equal(t,distance1,distance2, "distance in both direction should be same")
}


func TestValidPosition(t *testing.T) {
	cases := []validPositiontestCase{
		{models.Position{Latitude: 0, Longitude: 0,IsValid: false}, false},
		{models.Position{Latitude: 1.03256, Longitude: 102.36589,IsValid: true}, true},
		{models.Position{Latitude: -50.421352, Longitude: 102.36589,IsValid: true}, true},
		{models.Position{Latitude: 23.421352, Longitude: -102.36589,IsValid: true}, true},
		{models.Position{Latitude: -12.36523, Longitude: -104.24563,IsValid: true}, true},
		{models.Position{Latitude: -90, Longitude: -180,IsValid: true}, true},
		{models.Position{Latitude: 93.85612, Longitude: 102.36589,IsValid: false}, false},
		{models.Position{Latitude: 78.325, Longitude: -190.36981,IsValid: false}, false},

	}

	for _, tc := range cases {
		got := checkValidPosition(tc.position)
		assert.Equal(t, tc.expected, got, "test case failed for check valid position")
	}
}
