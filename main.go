package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
	"github.com/mursalinsk-qi/gps-distance-calculator/cache/redis"
	"github.com/mursalinsk-qi/gps-distance-calculator/calculator"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
	log "github.com/sirupsen/logrus"
)


func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}
type DeviceData struct {
	Latitude, Longitude, Distance ,Speed ,Count float64
}

var redisClient *redis.Client
var ctx=context.Background()
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
	redisInstance := cache.RedisInstance{Ctx:ctx}
	err := redisInstance.ConnectToRedis(host, port, password, db)
	if err != nil {
		log.WithFields(log.Fields{
			"message": "error in redis connection, returning from main function",
		}).Error(err)
		return
	}
	redisClient=redisInstance.RedisClient
	redisInstance.Subscribe(channel,CalculateTrackingData)

}

func CalculateTrackingData(message *redis.Message) {
	tracking_data := models.TrackingData{}
	if err := json.Unmarshal([]byte(message.Payload), &tracking_data); err != nil {
		log.Error(err)
	}
	current_latitude := tracking_data.GPS.Position.Latitude
	current_longitude := tracking_data.GPS.Position.Longitude
	current_speed:=tracking_data.Velocity.Speed
	redisHashKey := "gps_distance_calculator"
	redisHashField := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
	previous_redis_data := DeviceData{}
	isExists,err:= redisClient.HExists(ctx, redisHashKey, redisHashField).Result()
	if err!=nil{
		log.Error(err)
	}
	if isExists{
		val, err := redisClient.HGet(ctx, redisHashKey, redisHashField).Result()
		if err != nil {
			log.Error(err)
		}
		err2 := json.Unmarshal([]byte(val), &previous_redis_data)
		if err2 != nil {
			log.Error(err2)
		}
	} else {
		previous_redis_data.Latitude = current_latitude
		previous_redis_data.Longitude = current_longitude
		StoreValuesInRedis(&previous_redis_data,redisHashKey, redisHashField)
	}
	previous_latitude := previous_redis_data.Latitude
	previous_longitude := previous_redis_data.Longitude
	// Distance calculation
	distance := calculator.CalculateDistanceInMeter(current_latitude, current_longitude, previous_latitude, previous_longitude)
	// Updating Values in Redis cache
	previous_redis_data.Speed=current_speed+previous_redis_data.Speed
	previous_redis_data.Distance = distance + previous_redis_data.Distance
	previous_redis_data.Count=1+previous_redis_data.Count
	previous_redis_data.Latitude = current_latitude
	previous_redis_data.Longitude = current_longitude
	// Avarage Speed Calculation
	avgarageSpeed:=calculator.CalculateAvarageSpeed(previous_redis_data.Speed,previous_redis_data.Count)

	// Storing new values in Redis cache
	StoreValuesInRedis(&previous_redis_data,redisHashKey, redisHashField)
	log.Infof("id: %s,[%f,%f]-[%f,%f],dist:%f total dist: %f meter,speed %f,total %f,avarage %f", tracking_data.SourceId, previous_latitude, previous_longitude, current_latitude, current_longitude, distance, previous_redis_data.Distance,current_speed,previous_redis_data.Speed,avgarageSpeed)
	
}

func StoreValuesInRedis(previous_redis_data *DeviceData, redisHashKey, redisHashField string) {
	
	
	jsonValue, err := json.Marshal(previous_redis_data)
	if err != nil {
		log.Error(err)
	} else {
		redisClient.HSet(ctx, redisHashKey, redisHashField, jsonValue)
	}
}


