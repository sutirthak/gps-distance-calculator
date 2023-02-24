package cache

import (

	"context"

	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	RedisClient *redis.Client
	Ctx         context.Context
}
type RedisServices interface {
	ConnectToRedis(string,string,string,int) error
	Subscribe(string, func(*redis.Message))
	CalculateTrackingData(*redis.Message)
	StoreRedisData(*CoordinateData, float64, float64, float64, string, string)
}
