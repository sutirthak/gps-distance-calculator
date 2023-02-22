package main

import (
	"github.com/mursalinsk-qi/gps-distance-calculator/cache"
)

func init() {
	cache.ConnectToRedis()
}
func main() {
	cache.Subscriber()
}
