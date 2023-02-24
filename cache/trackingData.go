package cache

import (
	"encoding/json"
	"fmt"

	"github.com/mursalinsk-qi/gps-distance-calculator/calculator"
	"github.com/mursalinsk-qi/gps-distance-calculator/models"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)
type CoordinateData struct {
	Latitude, Longitude, Distance float64
}
func (r *RedisInstance) CalculateTrackingData(message *redis.Message) {
	tracking_data := models.TrackingData{}
	if err := json.Unmarshal([]byte(message.Payload), &tracking_data); err != nil {
		log.Error(err)
	}
	current_latitude := tracking_data.GPS.Position.Latitude
	current_longitude := tracking_data.GPS.Position.Longitude
	redisHashKey := "gps_distance_calculator"
	redisHashField := fmt.Sprintf("source_id:%s", tracking_data.SourceId)
	previous_redis_data := CoordinateData{}
	isExists, err := r.RedisClient.HExists(r.Ctx, redisHashKey, redisHashField).Result()
	if err != nil {
		log.Error(err)
	}
	if isExists {
		val, err := r.RedisClient.HGet(r.Ctx, redisHashKey, redisHashField).Result()
		if err != nil {
			log.Error(err)
		}
		err2 := json.Unmarshal([]byte(val), &previous_redis_data)
		if err2 != nil {
			log.Error(err2)
		}
	} else {
		r.StoreRedisData(&previous_redis_data, current_latitude, current_longitude, 0, redisHashKey, redisHashField)
	}
	previous_latitude := previous_redis_data.Latitude
	previous_longitude := previous_redis_data.Longitude
	distance := calculator.CalculateDistanceInMeter(current_latitude, current_longitude, previous_latitude, previous_longitude)
	r.StoreRedisData(&previous_redis_data, current_latitude, current_longitude, distance, redisHashKey, redisHashField)
	log.Infof("device id: %s, [%f,%f]-[%f-%f] , distance:%f total distance: %f meter",tracking_data.SourceId,previous_latitude,previous_longitude,current_latitude,current_longitude,distance,previous_redis_data.Distance)
}

func (r *RedisInstance)StoreRedisData(previous_redis_data *CoordinateData,current_latitude,current_longitude, distance float64,redisHashKey,redisHashField string){
	previous_redis_data.Latitude = current_latitude
	previous_redis_data.Longitude = current_longitude
	previous_redis_data.Distance = distance + previous_redis_data.Distance
	jsonValue, err := json.Marshal(previous_redis_data)
	if err != nil {
		log.Error(err)
	}else{
		r.RedisClient.HSet(r.Ctx, redisHashKey, redisHashField, jsonValue)
	}
}
