package models
type Coordinate struct {
	Latitude, Longitude float64
}
type DeviceData struct {
	Coordinate Coordinate
	Distance ,Speed ,Count float64
}