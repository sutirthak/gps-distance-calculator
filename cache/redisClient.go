package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var ctx = context.Background()
func ConnectToRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "fmdp-staging.ddnsfree.com:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong)
	if err != nil {
		panic(err)
		
	}else{
		fmt.Println("Connected to redis client")
	}
}

func GetRedisInstance() *redis.Client {
	return client
}
