package calculator

func CalculateAvarageSpeed(totalSpeed,count float64) float64 {
	if count==0{
		return 0
	}
	avarageSpeed:=totalSpeed/count
	return avarageSpeed
}