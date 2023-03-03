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
		log.Error(err)
	}
	current_position := models.Position{Latitude: tracking_data.GPS.Position.Latitude, Longitude: tracking_data.GPS.Position.Longitude}
	if current_position.Latitude == 0 && current_position.Longitude == 0 {
		return
	}
	current_speed := tracking_data.Velocity.Speed
	redisHashKey := "gps_distance_calculator"
	redisHashField := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
	previous_device_data := models.Trip{}
	isExists, err := redisClient.HExists(ctx, redisHashKey, redisHashField).Result()
	if err != nil {
		log.Error(err)
	}
	if isExists {
		val, err := redisClient.HGet(ctx, redisHashKey, redisHashField).Result()
		if err != nil {
			log.Error(err)
		}
		err2 := json.Unmarshal([]byte(val), &previous_device_data)
		if err2 != nil {
			log.Error(err2)
		}
	} else {
		previous_device_data.PrevPosition.Latitude = current_position.Latitude
		previous_device_data.PrevPosition.Longitude = current_position.Longitude
		StoreValuesInRedis(&previous_device_data, redisHashKey, redisHashField)
	}
	previous_position := models.Position{Latitude: previous_device_data.PrevPosition.Latitude, Longitude: previous_device_data.PrevPosition.Longitude}
	// Distance and speed calculation
	distance := calculator.CalculateDistanceInMeter(current_position, previous_position)
	avgarageSpeed := calculator.CalculateAvarageSpeed(previous_device_data, current_speed)
	// Updating Values in Redis cache
	previous_device_data.AvgSpeed = avgarageSpeed
	previous_device_data.Distance = distance + previous_device_data.Distance
	previous_device_data.DataCount = 1 + previous_device_data.DataCount
	previous_device_data.PrevPosition.Latitude = current_position.Latitude
	previous_device_data.PrevPosition.Longitude = current_position.Longitude
	// Avarage Speed Calculation

	// Storing new values in Redis cache
	StoreValuesInRedis(&previous_device_data, redisHashKey, redisHashField)
	log.Infof("id: %s,[%f,%f]-[%f,%f],dist:%f total dist: %f meter,speed %f, avarage %f", tracking_data.SourceId, previous_position.Latitude, previous_position.Longitude, current_position.Latitude, current_position.Longitude, distance, previous_device_data.Distance, current_speed, previous_device_data.AvgSpeed)

}

func StoreValuesInRedis(previous_device_data *models.Trip, redisHashKey, redisHashField string) {

	jsonValue, err := json.Marshal(previous_device_data)
	if err != nil {
		log.Error(err)
	} else {
		redisClient.HSet(ctx, redisHashKey, redisHashField, jsonValue)
	}
}
