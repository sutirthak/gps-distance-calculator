package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	cache "github.com/sutirthak/gps-distance-calculator/cache/redis"
	"github.com/sutirthak/gps-distance-calculator/models"
)
type ResponseMessage struct{
	Message string `json:"message"`
}
func GetTripData(c echo.Context) error {
	redisClient := cache.GetRedisClient()
	ctx := context.Background()
	id := c.Param("sourceid")
	redisHashKey := os.Getenv("FMDP_REDIS_HASHKEY")
	redisHashField := fmt.Sprintf("source_id:%s", id)
	val, err := redisClient.HGet(ctx, redisHashKey, redisHashField).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, ResponseMessage{"source id not found"})
	}
	trip := models.Trip{}
	err = json.Unmarshal([]byte(val), &trip)
	if err != nil {
		log.Error("invalid JSON for TripData")
		return c.JSON(http.StatusBadRequest, ResponseMessage{"something went wrong"})
	}
	return c.JSON(http.StatusOK, trip)
}
