package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	cache "github.com/sutirthak/gps-distance-calculator/cache/redis"
	"github.com/sutirthak/gps-distance-calculator/calculator"
	"github.com/sutirthak/gps-distance-calculator/models"
)
type Status string
const (
	Running  Status = "In progress"
	Finished Status = "Completed"
)
type ResponseMessage struct {
	Message string `json:"message"`
}
type ResponseData struct {
	PerPage    int                    `json:"per_page"`
	Page       int                    `json:"page"`
	TotalCount int                    `json:"total_count"`
	Data       *[]models.TrackingData `json:"data"`
}

type HistoricCalculationResult struct {
	Token  string       `json:"token"`
	Status Status       `json:"status"`
	Trip *models.Trip `json:"result"`
}

func getTripFromRedis(device_id, redisHashKey string, RedisInstance cache.RedisInstance) (models.Trip, error) {
	trip:=models.Trip{}
	value, err := RedisInstance.RedisClient.HGet(RedisInstance.Ctx, redisHashKey, device_id).Result()
	if err != nil {
		return trip, err
	}
	
	err = json.Unmarshal([]byte(value), &trip)
	if err != nil {
		return trip, err
	}
	return trip, nil
}

// This api is used to get current trip distance
func GetCurrentTrip(c echo.Context) error {
	cc := c.(*models.Context)
	device_id := c.Param("sourceid")
	redisHashKey := os.Getenv("FMDP_REDIS_HASHKEY_CURRENT")
	trip, err := getTripFromRedis(device_id, redisHashKey, cc.RedisInstance)
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, ResponseMessage{"device id not found"})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseMessage{"something went wrong"})
	}
	return c.JSON(http.StatusOK, trip)
}

/*
	This function is used to get trip distance on a historical data based on start and end time for a device
	-> If start and end time not given return an error
	-> Create an unique id (this will be used in future to show the result) and return it
	-> create a post request for calculation in background

*/


func GetHistoricTrip(c echo.Context) error {
	deviceid := c.Param("sourceid")
	startTime := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")
	if startTime == "" || endTime == "" {
		return c.JSON(http.StatusBadRequest, ResponseMessage{"Please provide start and end time"})
	}
	data, err :=getDevices(deviceid, startTime, endTime)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, ResponseMessage{"internal server error"})
	}
	
	token := uuid.New().String()
	go sendPostRequest(deviceid,token,data)
	return c.JSON(http.StatusAccepted, token)
}

// This function is used to get data from FMDP Server
func getDevices(deviceid, startTime, endTime string) (ResponseData, error) {
	host := os.Getenv("FMDP_SERVER_HOST")
	url := fmt.Sprintf("%s/formatted/devices/%s?start_time=%s&end_time=%s", host, deviceid,startTime,endTime)
	data := ResponseData{}
	response, err := http.Get(url)
	if err != nil {
		return data, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}
	
	json.Unmarshal(body, &data)
	return data, nil
}

// This function is used to send post request
func sendPostRequest(deviceid,token string,data ResponseData){
	jsonValue, _ := json.Marshal(data)
	host:=os.Getenv("LOCAL_HOST")
	url := fmt.Sprintf("%s/trip/%s?token=%s", host,deviceid, token)
	http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
}
func PostHistoricTrip(c echo.Context) error {
	cc := c.(*models.Context)
	token := c.QueryParam("token")
	deviceid := c.Param("sourceid")
	redisHashkey := os.Getenv("FMDP_REDIS_HASHKEY_HISTORIC")
	result := HistoricCalculationResult{Token: token, Status: Running, Trip: nil}
	storeValuesInRedis(result, &cc.RedisInstance,redisHashkey, token)
	data := new(ResponseData)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	redisHashField:=fmt.Sprintf("%s:%s",deviceid,token)
	if len(*data.Data)>0{
		for _, value := range *data.Data {
			CalculateTrackingData(value, &cc.RedisInstance, redisHashkey,redisHashField)
		}
		result.Trip=&models.Trip{}
		*result.Trip,_=getTripFromRedis(redisHashField, redisHashkey, cc.RedisInstance)
	}
	result.Status =Finished
	storeValuesInRedis(result, &cc.RedisInstance, redisHashkey, token)
	cc.RedisInstance.RedisClient.HDel(cc.RedisInstance.Ctx,redisHashkey,redisHashField)
	return c.JSON(http.StatusOK, "Ok")
}

// getting token status

func GetStatus(c echo.Context) error{
	cc := c.(*models.Context)
	token := c.Param("token")
	redisHashKey := os.Getenv("FMDP_REDIS_HASHKEY_HISTORIC")
	value, err := cc.RedisInstance.RedisClient.HGet(cc.RedisInstance.Ctx, redisHashKey,token).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, ResponseMessage{"token not found"})
	}
	result:=HistoricCalculationResult{}
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseMessage{"something went wrong"})
	}
	return c.JSON(http.StatusOK, result)

}

// ----- Trip Calculation----------------------------------------------------

func CalculateTrackingData(tracking_data models.TrackingData, r *cache.RedisInstance, redisHashKey,redisHashField string) {
	if tracking_data.GPS == nil {
		return
	}
	current_position := tracking_data.GPS.Position
	previous_position := &models.Position{}
	current_speed := 0.0
	if tracking_data.Velocity != nil {
		current_speed = tracking_data.Velocity.Speed
	}
	trip, err := getTripFromRedis(redisHashField, redisHashKey, *r)
	if err == redis.Nil {
		previous_position = current_position
	} else {
		previous_position = trip.PrevPosition
	}
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
	storeValuesInRedis(trip, r, redisHashKey, redisHashField)
	// log.Infof("id: %s,[%f,%f]-[%f,%f],dist:%f total dist: %f meter,speed %f, avarage %f", tracking_data.SourceId, previous_position.Latitude, previous_position.Longitude, current_position.Latitude, current_position.Longitude, distance, trip.Distance, current_speed, trip.AvgSpeed)
}

func storeValuesInRedis(value interface{}, r *cache.RedisInstance, redisHashKey, redisHashField string) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Error(err)
	} else {
		r.RedisClient.HSet(r.Ctx, redisHashKey, redisHashField, jsonValue)
	}
}
