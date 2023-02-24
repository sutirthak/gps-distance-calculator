package main

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mursalinsk-qi/gps-distance-calculator/cache"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}

func main() {
	// log formatting
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	})
	// Getting data from .env file
	host := os.Getenv("FMDP_HOST")
	port := os.Getenv("FMDP_PORT")
	password := os.Getenv("FMDP_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("FMDP_DB"))
	channel := os.Getenv("FMDP_CHANNEL")
	// Redis Connection
	redisInstance:=cache.RedisInstance{Ctx: context.Background()}
	err:=redisInstance.ConnectToRedis(host, port, password, db)
	if err != nil {
		log.WithFields(log.Fields{
			"message": "error in redis connection, returning from main function",
		}).Error(err)
		return
	}
	redisInstance.Subscribe(channel, redisInstance.CalculateTrackingData)

}
