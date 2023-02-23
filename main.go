package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mursalinsk-qi/gps-distance-calculator/cache"
	"github.com/mursalinsk-qi/gps-distance-calculator/calculator"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
	"github.com/redis/go-redis/v9"
)
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	// Getting data from .env file
	host := os.Getenv("FMDF_HOST")
	port := os.Getenv("FMDF_PORT")
	password := os.Getenv("FMDF_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("FMDF_DB"))
	channel := os.Getenv("FMDF_CHANNEL")

	// Redis Connection
	cache.ConnectToRedis(host, port, password, db)
	redisClient = cache.GetRedisInstance()
	cache.Subscriber(redisClient, ctx, channel, CalculateTrackingData)

}

func CalculateTrackingData(message *redis.Message) {
	tracking_data := models.TrackingData{}
	if err := json.Unmarshal([]byte(message.Payload), &tracking_data); err != nil {
		panic(err)
	}
	current_latitude:=tracking_data.GPS.Position.Latitude
	current_longitude:=tracking_data.GPS.Position.Longitude
	redisHashKey:="tracking_data_distance"
	redisHashField := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
	redis_data := models.RedisData{}
	isExists, err := redisClient.HExists(ctx, redisHashKey, redisHashField).Result()
	if err != nil {
		panic(err)
	}
	if isExists {
		val, err := redisClient.HGet(ctx,redisHashKey, redisHashField).Result()
		if err != nil {
			panic(err)
		}
		err2 := json.Unmarshal([]byte(val), &redis_data)
		if err2 != nil {
			panic(err2)
		}
	} else {
		StoreRedisData(&redis_data,current_latitude,current_longitude,0,redisHashKey,redisHashField)
	}
	previous_latitude:=redis_data.Latitude
	previous_longitude:=redis_data.Longitude
	distance := calculator.CalculateDistanceInMeter(current_latitude,current_longitude,previous_latitude,previous_longitude)
	StoreRedisData(&redis_data,current_latitude,current_longitude,distance,redisHashKey,redisHashField)
	fmt.Printf("%v: ",tracking_data.SourceId)
	fmt.Printf("prev loc(%v,%v) curr loc (%v,%v)",previous_latitude,previous_longitude,current_latitude,current_longitude)
	fmt.Printf(" distance : %v , total distance covered %v meter\n", distance, redis_data.Distance)
}

func StoreRedisData(redis_data *models.RedisData,current_latitude,current_longitude, distance float64,redisHashKey,redisHashField string){
	redis_data.Latitude = current_latitude
	redis_data.Longitude = current_longitude
	redis_data.Distance = distance + redis_data.Distance
	jsonValue, err := json.Marshal(redis_data)
	if err != nil {
		panic(err)
	}
	redisClient.HSet(ctx, redisHashKey, redisHashField, jsonValue)
}
