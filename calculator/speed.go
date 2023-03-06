package calculator

import "github.com/mursalinsk-qi/gps-distance-calculator/models"

func CalculateAvarageSpeed(tripData models.Trip, current_speed float64) float64 {
	totalSpeed:=tripData.AvgSpeed*float64(tripData.DataCount)+ current_speed
	avarageSpeed:=totalSpeed/ float64(tripData.DataCount+1)
	return avarageSpeed
}