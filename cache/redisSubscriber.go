package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mursalinsk-qi/gps-distance-calculator/calculator"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
)
type RedisData struct {
	Latitude, Longitude, Distance float64
}

func Subscriber() {
	redisClient := GetRedisInstance()
	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, "adapter_output_topic")
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		tracking_data := models.TrackingData{}
		if err := json.Unmarshal([]byte(msg.Payload), &tracking_data); err != nil {
			panic(err)
		}
		key := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
		isExists, err := redisClient.HExists(ctx, "tracking_data_distance", key).Result()
		if err != nil {
			panic(err)
		}
		redis_data := RedisData{}
		if isExists {
			val, err := redisClient.HGet(ctx, "tracking_data_distance", key).Result()
			if err != nil {
				panic(err)
			}
			err2 := json.Unmarshal([]byte(val), &redis_data)
			if err2 != nil {
				panic(err2)
			}
		}else{
			redis_data.Latitude = tracking_data.GPS.Position.Latitude
			redis_data.Longitude=tracking_data.GPS.Position.Longitude
		}
		distance:=calculator.CalculateDistanceInKm(tracking_data.GPS.Position.Latitude,tracking_data.GPS.Position.Longitude,redis_data.Latitude,redis_data.Longitude)
		fmt.Printf("prev position (%v,%v) curr position (%v,%v)",redis_data.Latitude,redis_data.Longitude,tracking_data.GPS.Position.Latitude,tracking_data.GPS.Position.Longitude)
		redis_data.Latitude = tracking_data.GPS.Position.Latitude
		redis_data.Longitude=tracking_data.GPS.Position.Longitude
		redis_data.Distance=distance+redis_data.Distance
		fmt.Printf(" distance : %v , total distance covered %v\n",distance,redis_data.Distance)
		jsonValue, err := json.Marshal(redis_data)
		if err != nil {
			panic(err)
		}
		redisClient.HSet(ctx, "tracking_data_distance", key, jsonValue)
	}
}

