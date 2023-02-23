package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var ctx = context.Background()

func ConnectToRedis(host, port, password string, db int) {
	adress := fmt.Sprintf("%s:%s", host, port)
	client = redis.NewClient(&redis.Options{
		Addr:     adress,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong)
	if err != nil {
		panic(err)

	} else {
		fmt.Println("Connected to redis client")
	}
}

func GetRedisInstance() *redis.Client {
	return client
}
