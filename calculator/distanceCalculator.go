package calculator

import (
	"math"
)

func CalculateDistanceInKm(current_latitude, current_longitude, previous_latitude, previous_longitude float64) float64{
	radiousOfEarth := float64(6371)
	dlat:=degreeToRadious(previous_latitude-current_latitude)
	dlot:=degreeToRadious(previous_longitude-current_longitude)
	a:=math.Sin(dlat/2)*math.Sin(dlat/2)+math.Cos(degreeToRadious(current_latitude))*math.Cos(degreeToRadious(previous_latitude))*math.Sin(dlot/2)*math.Sin(dlot/2)
	c:=2 * math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	distance:=radiousOfEarth*c
	return distance

}

func degreeToRadious(degree float64) float64 {
	return degree * (math.Pi /180)
}