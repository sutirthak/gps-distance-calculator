package calculator

import (
	"math"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
)

func CalculateDistanceInMeter(startingPosition,endPosition models.Position) float64{
	if !checkValidPosition(startingPosition) || !checkValidPosition(endPosition){
		return 0
	}
	radiousOfEarth := float64(6371)
	dlat := degreeToRadious(startingPosition.Latitude - endPosition.Latitude)
	dlot := degreeToRadious(startingPosition.Longitude - endPosition.Longitude)
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(degreeToRadious(endPosition.Latitude))*math.Cos(degreeToRadious(startingPosition.Latitude))*math.Sin(dlot/2)*math.Sin(dlot/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := radiousOfEarth * c
	return distance * 1000

}

func degreeToRadious(degree float64) float64 {
	return degree * (math.Pi /180)
}

func checkValidPosition(position models.Position) bool{
	if position.Latitude>90 || position.Latitude<(-90) {
		return false
	}
	if position.Longitude>180 || position.Longitude<(-180) {
		return false
	}
	return true
}