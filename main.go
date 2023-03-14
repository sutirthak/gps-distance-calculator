package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	cache "github.com/sutirthak/gps-distance-calculator/cache/redis"
	"github.com/sutirthak/gps-distance-calculator/controller"
	"github.com/sutirthak/gps-distance-calculator/models"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}

var version string

type Version struct {
	Version string `json:"version"`
}

func main() {
	// log formatting
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	})
	// Getting data from .env file
	host := os.Getenv("FMDP_REDIS_HOST")
	redis_server_port := os.Getenv("FMDP_PORT")
	password := os.Getenv("FMDP_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("FMDP_DB"))
	channel := os.Getenv("FMDP_CHANNEL")
	eco_server_port := os.Getenv("ECO_SERVER_PORT")

	// Redis Connection
	redisInstance := cache.RedisInstance{Ctx: context.Background()}
	err := redisInstance.ConnectToRedis(host, redis_server_port, password, db)
	if err != nil {
		log.WithFields(log.Fields{
			"message": "error in redis connection, returning from main function",
		}).Error(err)
		return
	}
	go redisInstance.Subscribe(channel, receiveMessage)
	// Server Connection
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &models.Context{
				Context: c, 
				RedisInstance: redisInstance,
			}
			return next(cc)
		}
	})
	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusAccepted, Version{version})
	})
	e.GET("/trip/current/:sourceid", controller.GetCurrentTrip)
	e.GET("/trip/historic/:sourceid", controller.GetHistoricTrip)
	e.GET("/trip/:token/status", controller.GetStatus)
	e.POST("/trip/:sourceid", controller.PostHistoricTrip)
	e.Logger.Fatal(e.Start(eco_server_port))

}
func receiveMessage(message *redis.Message, r *cache.RedisInstance) {
	tracking_data := models.TrackingData{}
	if err := json.Unmarshal([]byte(message.Payload), &tracking_data); err != nil {
		log.WithFields(log.Fields{
			"message": "invalid JSON for TrackingData",
		}).Error(err)
		return
	}
	// Calculating current distance
	redisHashkey := os.Getenv("FMDP_REDIS_HASHKEY_CURRENT")
	controller.CalculateTrackingData(tracking_data, r, redisHashkey,tracking_data.SourceId)
}
