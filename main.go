package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	cache "github.com/mursalinsk-qi/gps-distance-calculator/cache/redis"
	"github.com/mursalinsk-qi/gps-distance-calculator/calculator"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}

var redisClient *redis.Client
var ctx = context.Background()

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
	redisInstance := cache.RedisInstance{Ctx: ctx}
	err := redisInstance.ConnectToRedis(host, port, password, db)
	if err != nil {
		log.WithFields(log.Fields{
			"message": "error in redis connection, returning from main function",
		}).Error(err)
		return
	}
	redisClient = redisInstance.RedisClient
	redisInstance.Subscribe(channel, CalculateTrackingData)

}

func CalculateTrackingData(message *redis.Message) {
	tracking_data := models.TrackingData{}
	if err := json.Unmarshal([]byte(message.Payload), &tracking_data); err != nil {
		log.WithFields(log.Fields{
			"message": "invalid JSON for TrackingData",
		}).Error(err)
		return
	}
	current_position := tracking_data.GPS.Position
	current_speed := tracking_data.Velocity.Speed
	redisHashKey := "gps_distance_calculator"
	redisHashField := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
	trip := models.Trip{}
	isExists, err := redisClient.HExists(ctx, redisHashKey, redisHashField).Result()
	if err != nil {
		log.Error(err)
		return
	}
	if isExists {
		val, err := redisClient.HGet(ctx, redisHashKey, redisHashField).Result()
		if err != nil {
			log.Error(err)
			return
		}
		err2 := json.Unmarshal([]byte(val), &trip)
		if err2 != nil {
			log.WithFields(log.Fields{
				"message": "invalid JSON for TripData",
			}).Error(err)
			return
		}
	} else {
		trip.PrevPosition = current_position
	}
	previous_position := trip.PrevPosition
	// Distance and speed calculation
	distance, err := calculator.CalculateDistanceInMeter(*current_position, *previous_position)
	if err != nil {
		log.Error(err)
		return
	}
	avgarageSpeed := calculator.CalculateAvarageSpeed(trip, current_speed)
	// Updating Values in Redis cache
	trip.AvgSpeed = avgarageSpeed
	trip.Distance = distance + trip.Distance
	trip.DataCount = 1 + trip.DataCount
	trip.PrevPosition = current_position
	// Storing new values in Redis cache
	StoreValuesInRedis(&trip, redisHashKey, redisHashField)
	log.Infof("id: %s,[%f,%f]-[%f,%f],dist:%f total dist: %f meter,speed %f, avarage %f", tracking_data.SourceId, previous_position.Latitude, previous_position.Longitude, current_position.Latitude, current_position.Longitude, distance, trip.Distance, current_speed, trip.AvgSpeed)

}

func StoreValuesInRedis(trip *models.Trip, redisHashKey, redisHashField string) {

	jsonValue, err := json.Marshal(trip)
	if err != nil {
		log.Error(err)
	} else {
		redisClient.HSet(ctx, redisHashKey, redisHashField, jsonValue)
	}
}
